package goo_task_queue

import (
	"errors"
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
		if task.MaxRetry == 0 {
			task.MaxRetry = 99 // 默认重试次数 99次
		}
		if task.Timeout == 0 {
			task.Timeout = 1800 // 默认超时时间 30分钟
		}

		pi.HMSet(p.taskInfoKey(task.Id), task.MapData())
		pi.ZAdd(p.TaskPendingKey, redis.Z{Score: task.Score(), Member: task.Id})
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
	return p.TaskQueue.log().WithTag("goo-task-queue-publish")
}
