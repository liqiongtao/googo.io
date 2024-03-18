# 服务端

- `Config` 配置信息
    - `ServiceName` 服务名称
    - `ServiceEndpoint` 服务开放地址
    - `Addr` 服务监听地址，必须明确 `ip:port`
- 信号监控
    - `kill -USR1` 开启 `pprof` 监控
    - `kill -USR2` 停止 `pprof` 监控，可以获取监控文件
    - `kill -1` 平滑重启
    - `kill -9` 退出应用程序，目前监控不到
    - `kill -QUIT` 退出应用程序
    - `ctrl + C` 退出应用程序
- 拦截器
    - `grpc.ChainUnaryInterceptor` 服务端单向拦截器（也叫"一元拦截器"）
        - `serverUnaryInterceptorLog()` 记录日志信息
        - `unaryServerInterceptorAuth()` 鉴权处理
            - `AuthFunc()` 自定义鉴权方法
    - `grpc.ChainStreamInterceptor` 服务端流式拦截器
        - 同上
- `Register2ETCD()` 服务注册到etcd
    - `key` 格式: <service>/<leaseId>
    - `value` 格式: json格式，示例 `{"Op":0,"Addr":"127.0.0.1:19001","Metadata":null}`

# 客户端

- `Dial()` 链接 grpc
- `DialWithEtcd()` 通过 etcd 链接 grpc
    - `target` 格式 `etcd:///<service>`
- 返回拦截
  - 

# grpc 版本要求

最高使用 `v1.52.3`, 否则导致 `Target.Endpoint` 无效，因为高版本把属性定义为了方法
