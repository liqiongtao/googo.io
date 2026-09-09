package goo_task_queue

type TaskQueueList struct {
	*TaskQueue
}

// 待执行的任务Id集合
func (q *TaskQueueList) PendingTasks() []*Task {
	return q.getTasks(q.TaskPendingKey)
}

// 正在执行的任务Id集合
func (q *TaskQueueList) ProcessingTasks() []*Task {
	return q.getTasks(q.TaskProcessingKey)
}

// 执行失败的任务Id集合
func (q *TaskQueueList) FailTasks() []*Task {
	return q.getTasks(q.TaskFailKey)
}

// 执行失败的任务Id集合
func (q *TaskQueueList) getTasks(key string) []*Task {
	taskIds := q.r.ZRange(key, 0, -1).Val()

	var tasks []*Task
	for _, taskId := range taskIds {
		if taskId == "" {
			continue
		}
		exists, err := q.TaskQueueTasks.taskExists(taskId)
		if err != nil || !exists {
			continue
		}
		task, err := getTaskByCache(q.r, q.TaskQueueTasks.taskInfoKey(taskId))
		if err != nil || task.Id == "" {
			continue
		}
		tasks = append(tasks, task)
	}

	return tasks
}
