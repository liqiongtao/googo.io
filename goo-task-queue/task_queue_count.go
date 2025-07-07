package goo_task_queue

type TaskQueueCount struct {
	*TaskQueue
}

// 待执行的任务数量
func (q *TaskQueueCount) PendingCount() int64 {
	return q.r.ZCard(q.TaskPendingKey).Val()
}

// 正在执行的任务数量
func (q *TaskQueueCount) ProcessingCount() int64 {
	return q.r.ZCard(q.TaskProcessingKey).Val()
}

// 执行失败的任务数量
func (q *TaskQueueCount) FailCount() int64 {
	return q.r.ZCard(q.TaskFailKey).Val()
}
