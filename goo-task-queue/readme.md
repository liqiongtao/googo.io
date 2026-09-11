# goo-task-queue

基于 Redis 的分布式任务队列，支持发布/订阅、优先级、延迟调度、重试、超时恢复与 Worker 心跳管理。

时间单位统一为**毫秒**（`Timeout`、`Ts`、pending/processing score）。

## 功能概览

| 能力 | 说明 |
|------|------|
| 任务发布 | 写入任务详情，加入待执行队列；同 ID 覆盖，并从失败队列移除 |
| 延迟调度 | `Ts` 设为未来毫秒时间戳即可延迟执行 |
| 任务消费 | 多 Worker 并发拉取执行，默认并发数 = `CPU核数 * 2` |
| 优先级 | `HighPriority=1` 在同一可执行时刻优先调度 |
| 成功清理 | 执行成功后删除任务详情及各队列中的记录 |
| 失败重试 | 失败后重试次数 +1，重新进入待执行队列；可用 `RetryAfter` 由业务指定延迟 |
| 失败落盘 | 达到 `MaxRetry` 后转入失败队列 |
| 超时恢复 | Leader 扫描执行中超时任务，计入重试后重新入队（默认立即） |
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
    Timeout      int64  // 超时毫秒数（默认 2 小时）
    Ts           int64  // 调度时间（毫秒时间戳）；0 表示立即；未来时间表示延迟
    Generation   int64  // 执行代数（框架内部，业务勿写）
}
```

## pending score 规则

```text
score = Ts           // HighPriority == 1
score = Ts + 0.5     // 普通任务
```

拉取：`ZRANGEBYSCORE pending 0 (nowMs+0.5)`，因此：

- `Ts > now` → 未到期，不会被捞起（延迟执行）
- 同一时刻：高优排在普通之前

## Redis 数据结构

| Key | 类型 | 说明 |
|-----|------|------|
| `tq:task:info:{id}` | Hash | 任务详情，TTL 48 小时 |
| `tq:task:pending` | ZSet | 待执行队列，score 见上 |
| `tq:task:processing` | ZSet | 执行中队列，score=最近续租时间（毫秒），超时未续租则回收 |
| `tq:task:fail` | ZSet | 失败队列 |
| `tq:task:workers` | Hash | Worker 心跳，field=`IP:PID` |
| `tq:task:leader:lock` | String | Leader 选举锁，TTL 10 秒 |

支持 Redis `Prefix`，会自动加到所有 Key 前。

## 核心流程

```
Publish
  → 写入 task:info
  → generation+1（使旧执行收尾失效）
  → ZADD pending（priorityScore(HighPriority, Ts)）
  → ZREM processing / fail（执行中重投会摘掉 processing）

Subscribe
  → 竞选 Leader / 上报心跳
  → Lua 原子：pending → processing，并 generation+1（本次执行代数）
  → 执行中按 Timeout/3 续租（刷新 processing score）
  → 执行 handler
      ├─ 成功 → 校验 generation 后删除任务
      ├─ 失败/超时且未超 MaxRetry → 校验 generation 后重试入队
      │     · 返回 RetryAfter(d, err) → next_run_at = now + d
      │     · 普通 error → next_run_at = now（立即）
      └─ 达到 MaxRetry → 校验 generation 后转入 fail
      （generation 不匹配则跳过收尾，避免旧执行踩踏新状态）

Leader
  → 续期选举锁
  → 扫描 processing：超过 Timeout 未续租则 Lua 原子 generation+1 并 requeueOrFail
  → 清理超过 20 秒无心跳的 Worker
```

## 使用示例

### 发布任务

```go
r, _ := gooredis.New(gooredis.Config{
    Addr:   "127.0.0.1:6379",
    Prefix: "myapp",
})

q := gootaskqueue.New(r)

err := q.Publish(&gootaskqueue.Task{
    Id:           "order-1001",
    Type:         "order.notify",
    Payload:      `{"order_id":1001}`,
    HighPriority: 1,
    MaxRetry:     5,
    Timeout:      int64(10 * time.Minute / time.Millisecond), // 毫秒
    // Ts: time.Now().Add(time.Minute).UnixMilli(), // 可选：延迟 1 分钟
})
```

### 消费任务（业务控制重试延迟）

```go
q.Subscribe(4, func(ctx context.Context, task *gootaskqueue.Task) error {
    if err := doWork(ctx, task); err != nil {
        // 限流等场景：2 分钟后再试（业务不感知 Redis/score）
        return gootaskqueue.RetryAfter(2*time.Minute, err)
    }
    return nil
})
```

- 直接 `return err`：立即重新可调度（`now`）
- `return RetryAfter(d, err)`：`now+d` 后才可调度
- `limit <= 0` 时使用默认并发数 `runtime.NumCPU() * 2`
- 可通过 `q.WithMaxMemoryPercent(85)` 调整内存保护阈值

### 查询状态

```go
q.PendingCount()     // 待执行数量
q.ProcessingCount()  // 执行中数量
q.FailCount()        // 失败数量

q.PendingCountByType("order")      // 指定类型待执行数量
q.ProcessingCountByType("order")   // 指定类型执行中数量
q.TypeCount("order")               // 指定类型：pending + processing
q.PendingCountGroupByType()        // 待执行按类型分组
q.ProcessingCountGroupByType()     // 执行中按类型分组
q.TypeCounts()                     // 全部类型：pending + processing

q.PendingTasks()     // 待执行任务列表
q.ProcessingTasks()  // 执行中任务列表
q.FailTasks()        // 失败任务列表
```

## 设计要点

1. **投递语义**：at-least-once；消费逻辑须可重入、幂等（超时回收 / 重投递可能造成重复执行）。
2. **generation 收尾校验**：每次抢占/回收/重投抬高代数；Worker 成功/失败/重试写 Redis 前校验，过期执行不能改队列状态。
3. **业务隔离**：重试延迟通过公开 API `RetryAfter` 表达；业务不直接操作 score / Redis。
4. **去重 / 覆盖**：相同 `Id` 再次 Publish 会覆盖任务详情、摘掉 processing，并抬高 generation。
5. **原子抢占**：通过 Lua 脚本完成 `pending → processing` 并分配 generation。
6. **失败与超时**：业务失败与执行超时均校验 `MaxRetry`；未超限则重试，达到上限进入 fail。
7. **超时兜底**：Worker 崩溃无法续租时，Leader 在超过 Timeout 未续租后用单段 Lua 原子 bump generation 并重新入队/转 fail，避免 bump 成功但入队失败的中间态。
8. **执行续租**：执行中每 `Timeout/3`（且必须 `< Timeout`）刷新 processing 租约；续租随任务 ctx 结束（含 Timeout），避免卡死任务被无限续租。
9. **Worker 身份**：`LocalIP:PID`，心跳约每 200–800ms 刷新一次；超过 20 秒视为失效并清理。
10. **时间单位**：任务相关字段与调度 score 均为毫秒。
11. **内存保护**：`MaxMemoryPercent` 限制并发拉新，避免内存型任务把机器打爆。
