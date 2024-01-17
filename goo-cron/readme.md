# 定时任务

- 阻塞模式

```
c := goo_cron.Default()

c.Second(func() {
    fmt.Println("--11---", time.Now().Format("15:04:05"))
    time.Sleep(3 * time.Second)
    fmt.Println("--12---", time.Now().Format("15:04:05"))
})
c.SecondX(3, func() {
    fmt.Println("--21---", time.Now().Format("15:04:05"))
    time.Sleep(2 * time.Second)
    fmt.Println("--22---", time.Now().Format("15:04:05"))
})

c.Run()
```

- 非阻塞模式

```
c := goo_cron.Default()

c.Second(func() {
    fmt.Println("--11---", time.Now().Format("15:04:05"))
    time.Sleep(3 * time.Second)
    fmt.Println("--12---", time.Now().Format("15:04:05"))
})
c.SecondX(3, func() {
    fmt.Println("--21---", time.Now().Format("15:04:05"))
    time.Sleep(2 * time.Second)
    fmt.Println("--22---", time.Now().Format("15:04:05"))
})

c.Start()

<-goo_context.Cancel().Done()
```