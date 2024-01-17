# 定时任务

- 阻塞模式

```
c := goo_cron.Default()

c.Second(func() {
    fmt.Println("--1---", time.Now().Format("15:04:05"))
    time.Sleep(3 * time.Second)
    fmt.Println("--2---", time.Now().Format("15:04:05"))
})

c.Run()
```

- 非阻塞模式

```
c := goo_cron.Default()

c.Second(func() {
    fmt.Println("--1---", time.Now().Format("15:04:05"))
    time.Sleep(3 * time.Second)
    fmt.Println("--2---", time.Now().Format("15:04:05"))
})

c.Start()

<-goo_context.Cancel().Done()
```

# shell 检查任务是否退出

说明:

- `kill -1` 终端断线，重新加载
- `kill -2` 中断，同`ctrl+c`
- `kill -3` 退出
- `kill -9` 强制终止

```
cd /opt
kill -3 `pidof /opt/s.s`
for n in {1..60}; do
    sleep 1s
    pid=`pidof /opt/s.s`
    if [ "\$pid" == "" ]; then
        nohup /opt/s.s >> logs/log.log 2>&1 &
        break
    fi
done
```