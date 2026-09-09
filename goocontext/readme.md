# goocontext

统一的 Context 框架：**进程根 + 请求派生**。HTTP / gRPC 用 `gracenet`，信号只在本包 Notify 一次。

## 进程根与信号

| 档位 | 信号 | 命令 | 框架默认 | 钩子 |
|------|------|------|----------|------|
| Restart | `SIGHUP` | `kill -1` | 不 cancel | `OnRestart` |
| Exit | `SIGTERM` / `SIGINT` / `SIGQUIT` | `kill` / `Ctrl+C` / `kill -3` | 钩子后 `cancel(Root)`（只一次） | `OnExit` |
| Action | `SIGUSR1` | `kill -USR1` | 无 | `OnSignal`（pprof） |

钩子均异步派发，信号循环不因长 `Shutdown` 卡住；Exit 用 `sync.Once`，重复 TERM/INT/QUIT 只生效一次。进入 Exit 后忽略 HUP/USR*。`cancel(Root)` 仍在 Exit 钩子跑完之后，因此 `<-Root().Done()` 会等优雅退出结束。

**启动顺序**：先 `OnExit`/`OnRestart`，再让其它组件调用 `Root()`（如 `etcd.New`、`RegisterSignal`）。若钩子尚未注册就收到退出信号，会直接 `cancel(Root)`，优雅退出钩子来不及执行。

平滑重启只 `StartProcess`，父进程继续服务；子进程 `NotifyParentExitAfter` 后父进程走 Exit 退出。handoff 依赖子进程通知；父进程侧对 `restarting` 有超时复位，避免子进程异常时永久无法再 HUP。

```go
goocontext.OnExit(func() { /* 平滑退出 */ })
goocontext.OnRestart(func() { /* StartProcess */ })
goocontext.NotifyParentExitAfter(300 * time.Millisecond)
<-goocontext.Root().Done()
```

## 请求派生

在途请求/消费不要挂在 `Root()` 上，否则进程退出会取消未完成业务：

```go
ctx := goocontext.WithGenerateTraceId(context.Background())
log := goocontext.Log(ctx) // WithField 为 copy-on-write，需 log = log.WithField(...)
```
