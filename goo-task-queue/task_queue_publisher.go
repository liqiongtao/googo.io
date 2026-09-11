package gootaskqueue

import (
	"errors"
	"fmt"
	"time"

	goolog "github.com/liqiongtao/googo.io/goo-log"
	"github.com/liqiongtao/googo.io/goo-redis"
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

	ctx := p.r.Context()
	pi := p.r.TxPipeline()

	for _, task := range tasks {
		if task == nil {
			return fmt.Errorf("任务为空")
		}
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
		infoKey := p.taskInfoKey(task.Id)

		pi.HMSet(ctx, infoKey, task.MapData())
		pi.HIncrBy(ctx, infoKey, "generation", 1) // 使旧执行收尾失效
		pi.Expire(ctx, infoKey, time.Duration(infoTTLSecUntil(task.Ts))*time.Second)
		pi.ZAdd(ctx, p.TaskPendingKey, gooredis.Z{Score: score, Member: task.Id})
		pi.ZRem(ctx, p.TaskProcessingKey, task.Id) // 执行中重投：摘掉 processing
		pi.ZRem(ctx, p.TaskFailKey, task.Id)
	}

	if _, err := pi.Exec(ctx); err != nil {
		p.log().WithField("tasks", tasks).ErrorF("发布任务失败: %s", err.Error())
		return err
	}

	p.log().WithField("tasks", tasks).Info("发布任务成功")

	return nil
}

func (p *TaskQueuePublisher) log() *goolog.Entry {
	return goolog.WithTag("goo-task-queue-publish", p.pid)
}
