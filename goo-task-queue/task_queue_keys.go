package goo_task_queue

import "fmt"

var (
	defTaskLeadLockKey   = "tq:task:lead:lock"  // 任务选举锁
	defTaskWorkersKey    = "tq:task:workers"    // 任务工作节点
	defTaskInfoKey       = "tq:task:info"       // 任务信息
	defTaskPendingKey    = "tq:task:pending"    // 待处理任务 score=排序时间戳 从小到大排序
	defTaskProcessingKey = "tq:task:processing" // 正在处理任务 score=当前时间戳 用于判断超时
	defTaskFailKey       = "tq:task:fail"       // 失败任务
)

type TaskQueueKeys struct {
	TaskLeadLockKey   string // 任务选举锁 string
	TaskWorkersKey    string // 任务工作节点 hash
	TaskInfoKey       string // 任务信息 hash
	TaskPendingKey    string // 待处理任务 zset
	TaskProcessingKey string // 正在处理任务 zset
	TaskFailKey       string // 失败任务 zset
}

func NewTaskQueueKeys() *TaskQueueKeys {
	keys := &TaskQueueKeys{
		TaskLeadLockKey:   defTaskLeadLockKey,
		TaskWorkersKey:    defTaskWorkersKey,
		TaskInfoKey:       defTaskInfoKey,
		TaskPendingKey:    defTaskPendingKey,
		TaskProcessingKey: defTaskProcessingKey,
		TaskFailKey:       defTaskFailKey,
	}
	return keys
}

func (k *TaskQueueKeys) WithPrefix(prefix string) *TaskQueueKeys {
	k.TaskLeadLockKey = fmt.Sprintf("%s:%s", prefix, defTaskLeadLockKey)
	k.TaskWorkersKey = fmt.Sprintf("%s:%s", prefix, defTaskWorkersKey)
	k.TaskInfoKey = fmt.Sprintf("%s:%s", prefix, defTaskInfoKey)
	k.TaskPendingKey = fmt.Sprintf("%s:%s", prefix, defTaskPendingKey)
	k.TaskProcessingKey = fmt.Sprintf("%s:%s", prefix, defTaskProcessingKey)
	k.TaskFailKey = fmt.Sprintf("%s:%s", prefix, defTaskFailKey)
	return k
}

func (k *TaskQueueKeys) taskInfoKey(taskId string) string {
	return fmt.Sprintf("%s:%s", k.TaskInfoKey, taskId)
}
