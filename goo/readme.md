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
		<-goo_context.Cancel().Done()
	})

	wg.Wait()
}
```

# mysql

```
func main() {
	goo_db.Init(goo_context.Cancel(), goo_db.Config{
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
		<-goo_context.Cancel().Done()
	})
	wg.Wait()
}
```

# redis

```
func main() {
	goo_redis.Init(goo_context.Cancel(), goo_redis.Config{
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
		<-goo_context.Cancel().Done()
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
var cfg = goo_etcd.Config{
    User: "test",
    Password: "123456",
    Endpoints: []string{"localhost:23791", "localhost:23792", "localhost:23793"},
}

func init() {
	goo_etcd.Init(cfg)
}

func main() {
	s := goo.NewGRPCServer(goo_grpc.Config{
		ENV:         "test",
		ServiceName: "lpro/grpc-user",
		Version:     "v100",
		Addr:        "127.0.0.1:10011",
	}).Register2Etcd(goo_etcd.CLI())

	pb_user_v1.RegisterGetterServer(s.Server, &Server{})

	s.Serve()
}
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
