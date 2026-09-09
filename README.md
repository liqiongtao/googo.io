# github.com/liqiongtao/googo.io
golang development framework

# goo

> 知识点
- go 协程里面必须要用 recover，否则会导致主进程退出
- io流被读取后，会自动close

> 日志
- http请求日志，自动记录
- 主动打印日志，在需要的地方，手动调用

# goocontext

进程根 Context + 请求派生 + 统一信号。详见 [goocontext/readme.md](goocontext/readme.md)。

- 平滑重启：`kill -1`
- 平滑退出：`kill` / `kill -3` / `Ctrl+C`
- pprof 切换：`kill -USR1`

```go
goocontext.OnExit(func() { /* 平滑退出 */ })
goocontext.OnRestart(func() { /* 平滑重启 */ })
<-goocontext.Root().Done()
```
