# goo-task-queue

基于 Redis 的分布式任务队列，支持发布/订阅、优先级、重试、超时恢复与 Worker 心跳管理。

## 功能概览

| 能力 | 说明 |
|------|------|
| 任务发布 | 写入任务详情，加入待执行队列；同 ID 覆盖，并从失败队列移除 |
| 任务消费 | 多 Worker 并发拉取执行，默认并发数 = `CPU核数 * 2` |
| 优先级 | `HighPriority=1` 优先调度 |
| 成功清理 | 执行成功后删除任务详情及各队列中的记录 |
| 失败重试 | 失败后重试次数 +1，重新进入待执行队列 |
| 失败落盘 | 达到 `MaxRetry` 后转入失败队列 |
| 超时恢复 | Leader 扫描执行中超时任务，重新放回待执行队列 |
| Leader 选举 | Worker 通过 Redis 锁竞选 Leader，负责超时恢复与僵尸 Worker 清理 |
| 内存保护 | 内存占用超过阈值时暂停拉取新任务（默认 90%） |
| 优雅退出 | 响应取消信号，等待在途任务结束后再退出 |
| 监控查询 | 提供待执行 / 执行中 / 失败数量与列表查询 |

## 任务模型

```go
type Task struct {
    Id           string // 任务唯一 ID（必填）
    Type         string // 任务类型（必填）
    Payload      string // 任务数据（必填）
    HighPriority int    // 优先级：1=高，0=低（默认）
    MaxRetry     int    // 最大重试次数（默认 99）
    RetryTimes   int    // 当前重试次数
    Timeout      int64  // 超时秒数（默认 1800，即 30 分钟）
    Ts           int64  // 任务时间戳（用于调度排序）
}
```

## Redis 数据结构

| Key | 类型 | 说明 |
|-----|------|------|
| `tq:task:info:{id}` | Hash | 任务详情，TTL 48 小时 |
| `tq:task:pending` | ZSet | 待执行队列，score 为调度时间；高优先级 score=1 |
| `tq:task:processing` | ZSet | 执行中队列，score 为开始时间，用于判断超时 |
| `tq:task:fail` | ZSet | 失败队列 |
| `tq:task:workers` | Hash | Worker 心跳，field=`IP:PID` |
| `tq:task:leader:lock` | String | Leader 选举锁，TTL 10 秒 |

支持 Redis `Prefix`，会自动加到所有 Key 前。

## 核心流程

```
Publish
  → 写入 task:info
  → ZADD pending（高优先级 score=1，否则用时间戳）
  → ZREM fail（若曾失败可重新投递）

Subscribe
  → 竞选 Leader / 上报心跳
  → Lua 原子：pending → processing
  → 执行 handler
      ├─ 成功 → 删除任务
      ├─ 失败/超时且未超 MaxRetry → retry_times+1，回到 pending
      └─ 达到 MaxRetry → 转入 fail

Leader
  → 续期选举锁
  → 扫描 processing：超时则回到 pending
  → 清理超过 20 秒无心跳的 Worker
```

## 使用示例

### 发布任务

```go
r, _ := goo_redis.New(goo_redis.Config{
    Addr:   "127.0.0.1:6379",
    Prefix: "myapp",
})

q := goo_task_queue.New(r)

err := q.Publish(&goo_task_queue.Task{
    Id:           "order-1001",
    Type:         "order.notify",
    Payload:      `{"order_id":1001}`,
    HighPriority: 1,
    MaxRetry:     5,
    Timeout:      600, // 秒
})
```

### 消费任务

```go
q.Subscribe(4, func(ctx context.Context, task *goo_task_queue.Task) error {
    // ctx 带超时与 trace-id
    // 返回 nil 表示成功；返回 error 触发重试
    return handle(task)
})
```

- `limit <= 0` 时使用默认并发数 `runtime.NumCPU() * 2`
- 可通过 `q.WithMaxMemoryPercent(85)` 调整内存保护阈值

### 查询状态

```go
q.PendingCount()     // 待执行数量
q.ProcessingCount()  // 执行中数量
q.FailCount()        // 失败数量

q.PendingTasks()     // 待执行任务列表
q.ProcessingTasks()  // 执行中任务列表
q.FailTasks()        // 失败任务列表
```

## 设计要点

1. **投递语义**：at-least-once；消费逻辑须可重入、幂等（超时回收 / 重投递可能造成重复执行）。
2. **去重 / 覆盖**：相同 `Id` 再次 Publish 会覆盖任务详情，并确保进入 pending。
3. **原子抢占**：通过 Lua 脚本完成 `pending → processing`，避免多 Worker 抢到同一任务。
4. **失败与超时**：业务失败与执行超时均校验 `MaxRetry`；未超限则重试，达到上限进入 fail。
5. **超时兜底**：即便进程崩溃，Leader 也会把超时任务重新入队。
6. **Worker 身份**：`LocalIP:PID`，心跳约每 200–800ms 刷新一次；超过 20 秒视为失效并清理。
7. **默认超时**：`Timeout` 未设置时统一为 1800 秒（30 分钟）。
