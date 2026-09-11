package gootaskqueue

// priorityScore 计算 pending 队列 score：
// - 整数部分 / 主体：下次可执行时间（毫秒时间戳），支持延迟到未来
// - 高优先级：score = nextRunAtMs
// - 普通优先级：score = nextRunAtMs + 0.5（同一时刻高优优先）
func priorityScore(highPriority int, nextRunAtMs int64) float64 {
	if highPriority == 1 {
		return float64(nextRunAtMs)
	}
	return float64(nextRunAtMs) + 0.5
}
