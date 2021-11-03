# goo-etcd

## set/get

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
	key := "/test/proj.com/grpc-user/v100"
	val := "hnatao"

	goo_etcd.Set(key, val, 0)

	goo_log.Debug(goo_etcd.Get(key))
}
```

## GetMap

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
	key := "/test/proj.com/grpc-user"

	key1 := "/test/proj.com/grpc-user/v100"
	val1 := "hnatao-1"

	key2 := "/test/proj.com/grpc-user/v101"
	val2 := "hnatao-2"

	goo_etcd.Set(key1, val1, 0)
	goo_etcd.Set(key2, val2, 0)

	goo_log.Debug(goo_etcd.GetMap(key))
}
```

## GetArray

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
	key := "/test/proj.com/grpc-user"

	key1 := "/test/proj.com/grpc-user/v100"
	val1 := "hnatao-1"

	key2 := "/test/proj.com/grpc-user/v101"
	val2 := "hnatao-2"

	goo_etcd.Set(key1, val1, 0)
	goo_etcd.Set(key2, val2, 0)

	goo_log.Debug(goo_etcd.GetArray(key))
}
```

## SetWithKeepAlive

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
	key := "/test/proj.com/grpc-user/v100"
	val := "hnatao-1"

	goo_etcd.SetWithKeepAlive(key, val, 15)

	goo_utils.AsyncFunc(func() {
		for {
			select {
			case <-goo_context.Cancel().Done():
				return
			default:
				goo_log.Debug(goo_etcd.Get(key))
				time.Sleep(time.Second)
			}
		}
	})

	<-goo_context.Timeout(20 * time.Second).Done()
}
```