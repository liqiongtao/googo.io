# goo-grpc

## 服务端

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
	s := goo_grpc.New(goo_grpc.Config{
		ENV:         "test",
		ServiceName: "lpro/grpc-user",
		Version:     "v100",
		Addr:        "127.0.0.1:10011",
	}).Register2Etcd(goo_etcd.CLI())

	pb_user_v1.RegisterGetterServer(s.Server, &Server{})

	s.Serve()
}
```

## 客户端

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
	serviceName := "/test/lproj/grpc-user/v100"

	cc, err := goo_grpc.DialWithEtcd(context.TODO(), serviceName, goo_etcd.CLI())
	if err != nil {
		goo_log.Error(err.Error())
	}

	// -------------------------------------
	// -------------------------------------

	cli := pb_user_v1.NewGetterClient(cc)

	for {
		select {
		case <-goo_context.Cancel().Done():
			return
		default:
			cli.Info(context.TODO(), &pb_user_v1.InfoParams{Id: 1})
			time.Sleep(time.Second)
		}
	}
}
```