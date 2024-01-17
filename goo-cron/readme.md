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

```
cd /opt
kill -9 `pidof /opt/s.s`
for n in {1..60}; do
    sleep 1s
    pid=`pidof /opt/s.s`
    if [ "\$pid" == "" ]; then
        nohup /opt/s.s >> logs/log.log 2>&1 &
        break
    fi
done
```