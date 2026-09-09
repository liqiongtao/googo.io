package goo_task_queue

type TaskQueueCount struct {
	*TaskQueue
}

// TypeQueueCount 某一任务类型在待执行 / 执行中队列中的数量
type TypeQueueCount struct {
	Pending    int64 `json:"pending"`
	Processing int64 `json:"processing"`
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

// PendingCountByType 指定类型的待执行数量
func (q *TaskQueueCount) PendingCountByType(taskType string) int64 {
	return q.countType(q.TaskPendingKey, taskType)
}

// ProcessingCountByType 指定类型的执行中数量
func (q *TaskQueueCount) ProcessingCountByType(taskType string) int64 {
	return q.countType(q.TaskProcessingKey, taskType)
}

// TypeCount 指定类型的待执行 + 执行中数量
func (q *TaskQueueCount) TypeCount(taskType string) TypeQueueCount {
	return TypeQueueCount{
		Pending:    q.countType(q.TaskPendingKey, taskType),
		Processing: q.countType(q.TaskProcessingKey, taskType),
	}
}

// PendingCountGroupByType 按类型统计待执行数量
func (q *TaskQueueCount) PendingCountGroupByType() map[string]int64 {
	return q.countGroupByType(q.TaskPendingKey)
}

// ProcessingCountGroupByType 按类型统计执行中数量
func (q *TaskQueueCount) ProcessingCountGroupByType() map[string]int64 {
	return q.countGroupByType(q.TaskProcessingKey)
}

// TypeCounts 按类型汇总待执行 + 执行中数量
func (q *TaskQueueCount) TypeCounts() map[string]TypeQueueCount {
	pending := q.countGroupByType(q.TaskPendingKey)
	processing := q.countGroupByType(q.TaskProcessingKey)

	out := make(map[string]TypeQueueCount, len(pending)+len(processing))
	for typ, n := range pending {
		c := out[typ]
		c.Pending = n
		out[typ] = c
	}
	for typ, n := range processing {
		c := out[typ]
		c.Processing = n
		out[typ] = c
	}
	return out
}

func (q *TaskQueueCount) countType(zsetKey, taskType string) int64 {
	if taskType == "" {
		return 0
	}
	var n int64
	for _, taskId := range q.r.ZRange(zsetKey, 0, -1).Val() {
		if taskId == "" || !q.TaskQueueTasks.taskExists(taskId) {
			continue
		}
		if q.r.HGet(q.taskInfoKey(taskId), "type").Val() == taskType {
			n++
		}
	}
	return n
}

func (q *TaskQueueCount) countGroupByType(zsetKey string) map[string]int64 {
	taskIds := q.r.ZRange(zsetKey, 0, -1).Val()
	counts := make(map[string]int64)
	for _, taskId := range taskIds {
		if taskId == "" || !q.TaskQueueTasks.taskExists(taskId) {
			continue
		}
		typ := q.r.HGet(q.taskInfoKey(taskId), "type").Val()
		if typ == "" {
			continue
		}
		counts[typ]++
	}
	return counts
}
