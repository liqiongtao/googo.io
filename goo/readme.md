# CancelContext

```
func main() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("-----begin 1----")
		<-goo.CancelContext().Done()
		fmt.Println("-----end 1----")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("-----begin 2----")
		<-goo.CancelContext().Done()
		fmt.Println("-----end 2----")
	}()

	wg.Wait()
}
```

# TimeoutContext

```
func main() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("-----begin 1----")
		<-goo.TimeoutContext(3 * time.Second).Done()
		fmt.Println("-----end 1----")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("-----begin 2----")
		<-goo.TimeoutContext(5 * time.Second).Done()
		fmt.Println("-----end 2----")
	}()

	wg.Wait()
}
```

# 定时任务

```
func main() {
	var wg sync.WaitGroup

	goo.Crond(1*time.Second, func() {
		fmt.Println(time.Now().Format("15:04:05"))
	})

	goo.Crond(3*time.Second, func() {
		fmt.Println(time.Now().Format("15:04:05"))
	})

	goo.CronMinute(func() {
		fmt.Println(time.Now().Format("15:04:05"))
	})

	goo.CronHour(func() {
		fmt.Println(time.Now().Format("15:04:05"))
	})

	goo.CronDay(func() {
		fmt.Println(time.Now().Format("15:04:05"))
	})

	wg.Add(1)
	goo_utils.AsyncFunc(func() {
		defer wg.Done()
		<-goo.CancelContext().Done()
	})

	wg.Wait()
}
```

# mysql

```
func main() {
	goo_db.Init(goo.CancelContext(), goo_db.Config{
		Name:   "",
		Driver: "mysql",
		Master:      "root:123456@tcp(192.168.1.100:3306)/ttxian",
		Slaves:      []string{"root:123456@tcp(192.168.1.100:3307)/ttxian"},
		LogModel:    true,
		MaxIdle:     10,
		MaxOpen:     100,
		AutoPing:    true,
		LogFilePath: "",
		LogFileName: "",
	})

	m := map[string]string{}
	exist, err := goo.DB().Table("s_user").Get(&m)
	if err != nil {
		goo_log.Fatal(err.Error())
	}
	if !exist {
		goo_log.Fatal("no data")
	}
	goo_log.Debug(m["account"])

	var wg sync.WaitGroup
	wg.Add(1)
	goo_utils.AsyncFunc(func() {
		defer wg.Done()
		<-goo.CancelContext().Done()
	})
	wg.Wait()
}
```

# redis

```
func main() {
	goo_redis.Init(goo.CancelContext(), goo_redis.Config{
		Name:     "",
		Addr:     "192.168.1.100:6379",
		Password: "123456",
		DB:       0,
		Prefix:   "tt",
		AutoPing: true,
	})

	err := goo.Redis().Set("name", "hnatao", 3*time.Second).Err()
	if err != nil {
		goo_log.Fatal(err.Error())
	}

	name := goo.Redis().Get("name").String()
	goo_log.Debug(name)

	var wg sync.WaitGroup
	wg.Add(1)
	goo_utils.AsyncFunc(func() {
		defer wg.Done()
		<-goo.CancelContext().Done()
	})
	wg.Wait()
}
```

# Token

```
func main() {
	tokenStr, err := goo.CreateToken("1111", 100)
	if err != nil {
		goo_log.Fatal(err.Error())
	}
	goo_log.Debug(tokenStr)

	token, err := goo.ParseToken(tokenStr, "1111")
	if err != nil {
		goo_log.Fatal(err.Error())
	}
	goo_log.Debug(token)
}
```

# grpc

```
```

# server

```
type User struct {
	Age int `form:"age"`
}

func (u User) DoHandle(ctx *goo.Context) *goo.Response {
	if err := ctx.ShouldBind(&u); err != nil {
		return goo.Error(5001, "参数错误", err.Error())
	}
	return goo.Success(u.Age)
}

func main() {
	s := goo.NewServer()

	s.GET("/", goo.Handler(User{}))

	s.Run(":8080")
}
```
