package goo_task_queue

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
)

type TaskQueuePublisher struct {
	*TaskQueue
}

func (p *TaskQueuePublisher) Publish(tasks ...*Task) error {
	if p.r == nil {
		err := errors.New("redis未初始化")
		p.log().Error(err)
		return err
	}

	pi := p.r.TxPipeline()

	for _, task := range tasks {
		if task.Id == "" {
			return fmt.Errorf("任务Id为空")
		}
		if task.Type == "" {
			return fmt.Errorf("任务Type为空")
		}

		if task.MaxRetry == 0 {
			task.MaxRetry = defaultMaxRetry
		}
		if task.Timeout == 0 {
			task.Timeout = defaultTaskTimeout
		}

		if task.Ts == 0 {
			task.Ts = time.Now().UnixMilli()
		}

		score := priorityScore(task.HighPriority, task.Ts)

		pi.HMSet(p.taskInfoKey(task.Id), task.MapData())
		pi.Expire(p.taskInfoKey(task.Id), 48*time.Hour)
		pi.ZAdd(p.TaskPendingKey, redis.Z{Score: score, Member: task.Id})
		pi.ZRem(p.TaskFailKey, task.Id)
	}

	if _, err := pi.Exec(); err != nil {
		p.log().WithField("tasks", tasks).ErrorF("发布任务失败: %s", err.Error())
		return err
	}

	p.log().WithField("tasks", tasks).Info("发布任务成功")

	return nil
}

func (p *TaskQueuePublisher) log() *goo_log.Entry {
	return goo_log.WithTag("goo-task-queue-publish", p.pid)
}
