# goocontext

统一的 Context 框架：**进程根 + 请求派生**。HTTP / gRPC 热重启使用 `cloudflare/tableflip`，信号只在本包 Notify 一次。

## 进程根与信号

| 档位 | 信号 | 命令 | 框架默认 | 钩子 |
|------|------|------|----------|------|
| Restart | `SIGHUP` | `kill -1` | 不 cancel | `OnRestart` → `upg.Upgrade()` |
| Exit | `SIGTERM` / `SIGINT` / `SIGQUIT` | `kill` / `Ctrl+C` / `kill -3` | 钩子后 `cancel(Root)`（只一次） | `OnExit` |
| Action | `SIGUSR1` | `kill -USR1` | 无 | `OnSignal`（pprof） |

钩子均异步派发，信号循环不因长 `Shutdown` 卡住；Exit 用 `sync.Once`，重复 TERM/INT/QUIT 只生效一次。进入 Exit 后忽略 HUP/USR*。`cancel(Root)` 仍在 Exit 钩子跑完之后，因此 `<-Root().Done()` 会等优雅退出结束。

**启动顺序**：先 `OnExit`/`OnRestart`，再让其它组件调用 `Root()`（如 `etcd.New`、`RegisterSignal`）。若钩子尚未注册就收到退出信号，会直接 `cancel(Root)`，优雅退出钩子来不及执行。

平滑重启由 **tableflip** 完成：子进程 `Ready()` 后父进程 `Exit()`，再走 Exit 钩子优雅退出。子进程初始化失败或超时时，父进程可继续服务。

```go
upg, _ := tableflip.New(tableflip.Options{})
defer upg.Stop()

goocontext.OnExit(func() { /* Shutdown / GracefulStop */ })
goocontext.OnRestart(func() { _ = upg.Upgrade() })

ln, _ := upg.Listen("tcp", addr)
go serve(ln)
_ = upg.Ready()

go func() {
	<-upg.Exit()
	_ = syscall.Kill(os.Getpid(), syscall.SIGTERM)
}()
<-goocontext.Root().Done()
```

## 请求派生

在途请求/消费不要挂在 `Root()` 上，否则进程退出会取消未完成业务：

```go
ctx := goocontext.WithGenerateTraceId(context.Background())
log := goocontext.Log(ctx) // WithField 为 copy-on-write，需 log = log.WithField(...)
```
