package goo_task_queue

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/liqiongtao/googo.io/goo-redis"
)

var errStaleGeneration = errors.New("stale task generation")

// withGeneration 仅当 Hash.generation == expect 时执行后续命令；返回是否匹配并执行成功
func (t *TaskQueueTasks) evalWithGeneration(luaBody string, keys []string, expectGen int64, args ...any) (bool, error) {
	script := `
local gen = redis.call('HGET', KEYS[1], 'generation')
if not gen then
	gen = '0'
end
if tostring(gen) ~= tostring(ARGV[1]) then
	return 0
end
` + luaBody + `
return 1
`
	allArgs := append([]any{strconv.FormatInt(expectGen, 10)}, args...)
	res, err := t.r.Eval(script, keys, allArgs...).Result()
	if err != nil {
		if errors.Is(err, goo_redis.ErrNil) {
			return false, nil
		}
		return false, err
	}
	switch v := res.(type) {
	case int64:
		return v == 1, nil
	case int:
		return v == 1, nil
	default:
		return false, fmt.Errorf("unexpected lua result: %T %v", res, res)
	}
}

// renewLease 续租：仅当 generation 匹配且仍在 processing 中时刷新 score（避免收尾后重新写回）
func (t *TaskQueueTasks) renewLease(task *Task) error {
	if task == nil || task.Id == "" {
		return nil
	}
	nowMs := time.Now().UnixMilli()
	ok, err := t.evalWithGeneration(`
if redis.call('ZSCORE', KEYS[2], ARGV[3]) == false then
	return 0
end
redis.call('ZADD', KEYS[2], ARGV[2], ARGV[3])
`, []string{
		t.taskInfoKey(task.Id),
		t.TaskProcessingKey,
	}, task.Generation, float64(nowMs), task.Id)
	if err != nil {
		t.log().WithTag("renewLease").Error(err)
		return err
	}
	if !ok {
		return errStaleGeneration
	}
	return nil
}

func (t *TaskQueueTasks) taskDel(task *Task) error {
	if task == nil || task.Id == "" {
		return nil
	}
	ok, err := t.evalWithGeneration(`
redis.call('DEL', KEYS[1])
redis.call('ZREM', KEYS[2], ARGV[2])
redis.call('ZREM', KEYS[3], ARGV[2])
redis.call('ZREM', KEYS[4], ARGV[2])
`, []string{
		t.taskInfoKey(task.Id),
		t.TaskPendingKey,
		t.TaskProcessingKey,
		t.TaskFailKey,
	}, task.Generation, task.Id)
	if err != nil {
		t.log().WithTag("taskDel").Error(err)
		return err
	}
	if !ok {
		t.log().WithTag("taskDel").WithField("task_id", task.Id).WithField("generation", task.Generation).
			Warn(errStaleGeneration.Error())
		return errStaleGeneration
	}
	return nil
}

// taskDelForce 无 generation 校验的强制删除（脏数据清理）
func (t *TaskQueueTasks) taskDelForce(taskId string) error {
	if taskId == "" {
		return nil
	}
	ctx := t.r.Context()
	pi := t.r.TxPipeline()
	pi.Del(ctx, t.taskInfoKey(taskId))
	pi.ZRem(ctx, t.TaskPendingKey, taskId)
	pi.ZRem(ctx, t.TaskProcessingKey, taskId)
	pi.ZRem(ctx, t.TaskFailKey, taskId)
	if _, err := pi.Exec(ctx); err != nil {
		t.log().WithTag("taskDelForce").Error(err)
		return err
	}
	return nil
}

func (t *TaskQueueTasks) requeue(task *Task, nextRunAtMs int64) error {
	if task == nil || task.Id == "" {
		return nil
	}
	if nextRunAtMs <= 0 {
		nextRunAtMs = time.Now().UnixMilli()
	}
	score := priorityScore(task.HighPriority, nextRunAtMs)

	ok, err := t.evalWithGeneration(`
redis.call('HINCRBY', KEYS[1], 'retry_times', 1)
redis.call('HSET', KEYS[1], 'ts', ARGV[3])
redis.call('ZREM', KEYS[2], ARGV[2])
redis.call('ZADD', KEYS[3], ARGV[4], ARGV[2])
`, []string{
		t.taskInfoKey(task.Id),
		t.TaskProcessingKey,
		t.TaskPendingKey,
	}, task.Generation, task.Id, strconv.FormatInt(nextRunAtMs, 10), score)
	if err != nil {
		t.log().WithTag("requeue").Error(err)
		return err
	}
	if !ok {
		t.log().WithTag("requeue").WithField("task_id", task.Id).WithField("generation", task.Generation).
			Warn(errStaleGeneration.Error())
		return errStaleGeneration
	}
	return nil
}

func (t *TaskQueueTasks) taskFail(task *Task) error {
	if task == nil || task.Id == "" {
		return nil
	}
	ok, err := t.evalWithGeneration(`
redis.call('ZREM', KEYS[2], ARGV[2])
redis.call('ZADD', KEYS[3], ARGV[3], ARGV[2])
`, []string{
		t.taskInfoKey(task.Id),
		t.TaskProcessingKey,
		t.TaskFailKey,
	}, task.Generation, task.Id, float64(time.Now().UnixMilli()))
	if err != nil {
		t.log().WithTag("taskFail").Error(err)
		return err
	}
	if !ok {
		t.log().WithTag("taskFail").WithField("task_id", task.Id).WithField("generation", task.Generation).
			Warn(errStaleGeneration.Error())
		return errStaleGeneration
	}
	return nil
}

// putBack 将已抢占但未执行的任务放回 pending，不增加 retry_times（进程退出等）
func (t *TaskQueueTasks) putBack(task *Task) error {
	if task == nil || task.Id == "" {
		return nil
	}
	nextRunAtMs := time.Now().UnixMilli()
	score := priorityScore(task.HighPriority, nextRunAtMs)

	ok, err := t.evalWithGeneration(`
redis.call('HSET', KEYS[1], 'ts', ARGV[3])
redis.call('ZREM', KEYS[2], ARGV[2])
redis.call('ZADD', KEYS[3], ARGV[4], ARGV[2])
`, []string{
		t.taskInfoKey(task.Id),
		t.TaskProcessingKey,
		t.TaskPendingKey,
	}, task.Generation, task.Id, strconv.FormatInt(nextRunAtMs, 10), score)
	if err != nil {
		t.log().WithTag("putBack").Error(err)
		return err
	}
	if !ok {
		t.log().WithTag("putBack").WithField("task_id", task.Id).WithField("generation", task.Generation).
			Warn(errStaleGeneration.Error())
		return errStaleGeneration
	}
	return nil
}

// requeueOrFail Leader 超时回收：原子完成
// 仍在 processing → generation+1 → 达 MaxRetry 则进 fail，否则 retry_times+1 回 pending。
// 已不在 processing / info 不存在视为幂等成功（可能已被 Worker 收尾）。
func (t *TaskQueueTasks) requeueOrFail(task *Task, nextRunAtMs int64) error {
	if task == nil || task.Id == "" {
		return nil
	}
	if nextRunAtMs <= 0 {
		nextRunAtMs = time.Now().UnixMilli()
	}
	score := priorityScore(task.HighPriority, nextRunAtMs)

	script := `
-- KEYS[1]=info KEYS[2]=processing KEYS[3]=pending KEYS[4]=fail
-- ARGV[1]=taskId ARGV[2]=nextRunAtMs ARGV[3]=pendingScore ARGV[4]=defaultMaxRetry
if redis.call('ZSCORE', KEYS[2], ARGV[1]) == false then
	return 0
end
if redis.call('EXISTS', KEYS[1]) == 0 then
	redis.call('ZREM', KEYS[2], ARGV[1])
	return 0
end

local maxRetry = tonumber(redis.call('HGET', KEYS[1], 'max_retry') or '0') or 0
if maxRetry == 0 then
	maxRetry = tonumber(ARGV[4]) or 0
end
local retryTimes = tonumber(redis.call('HGET', KEYS[1], 'retry_times') or '0') or 0

redis.call('HINCRBY', KEYS[1], 'generation', 1)

if maxRetry ~= 0 and retryTimes >= maxRetry then
	redis.call('ZREM', KEYS[2], ARGV[1])
	redis.call('ZADD', KEYS[4], ARGV[2], ARGV[1])
	return 2
end

redis.call('HINCRBY', KEYS[1], 'retry_times', 1)
redis.call('HSET', KEYS[1], 'ts', ARGV[2])
redis.call('ZREM', KEYS[2], ARGV[1])
redis.call('ZADD', KEYS[3], ARGV[3], ARGV[1])
return 1
`
	res, err := t.r.Eval(script, []string{
		t.taskInfoKey(task.Id),
		t.TaskProcessingKey,
		t.TaskPendingKey,
		t.TaskFailKey,
	}, task.Id, strconv.FormatInt(nextRunAtMs, 10), score, defaultMaxRetry).Result()
	if err != nil {
		if errors.Is(err, goo_redis.ErrNil) {
			return nil
		}
		t.log().WithTag("requeueOrFail").WithField("task_id", task.Id).Error(err)
		return err
	}
	switch v := res.(type) {
	case int64:
		if v == 0 || v == 1 || v == 2 {
			return nil
		}
	case int:
		if v == 0 || v == 1 || v == 2 {
			return nil
		}
	}
	err = fmt.Errorf("unexpected lua result: %T %v", res, res)
	t.log().WithTag("requeueOrFail").Error(err)
	return err
}
