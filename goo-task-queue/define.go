package goo_task_queue

const (
	TaskInfoKey       = "task:queue:info"       // 任务信息
	TaskPendingKey    = "task:queue:pending"    // 待处理任务
	TaskProcessingKey = "task:queue:processing" // 正在处理任务
	TaskFailKey       = "task:queue:fail"       // 失败任务
)
