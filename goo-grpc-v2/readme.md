# 服务端

- 运行环境
- 服务名称
- 服务监听IP端口
- 性能分析 pprof
- 系统日志
- hook
- 信号监控
    - `syscall.SIGUSR1` 将 `goroutine` 状况 `dump` 下载，进行分析使用
    - `syscall.SIGUSR2` 开启/关闭所有监控指标，自行控制 `profiling` 时长
    - `syscall.SIGTERM` 真正开启 `gracefulStop`，优雅关闭
        - 如果是注册在etcd上，优雅关闭之前，需要从etcd上删掉注册的key

# 客户端
