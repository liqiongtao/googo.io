package gooetcd

import (
	"errors"
	"sync"

	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	__client *Client
	__mu     sync.RWMutex
)

func Init(conf Config) (err error) {
	cli, err := New(conf)
	if err != nil {
		return err
	}
	__mu.Lock()
	__client = cli
	__mu.Unlock()

	// 绑定本次实例，避免重复 Init 后误关新 client / 泄漏旧 client 清理
	goo_utils.AsyncFunc(func() {
		<-cli.ctx.Done()
		cli.Close()
	})

	return
}

func Default() *Client {
	__mu.RLock()
	defer __mu.RUnlock()
	return __client
}

func requireClient() (*Client, error) {
	__mu.RLock()
	cli := __client
	__mu.RUnlock()
	if cli == nil {
		return nil, errors.New("etcd not initialized, call gooetcd.Init first")
	}
	return cli, nil
}

func Set(key, val string) (resp *clientv3.PutResponse, err error) {
	cli, err := requireClient()
	if err != nil {
		return nil, err
	}
	return cli.Set(key, val)
}

func SetWithPrevKV(key, val string) (resp *clientv3.PutResponse, err error) {
	cli, err := requireClient()
	if err != nil {
		return nil, err
	}
	return cli.SetWithPrevKV(key, val)
}

func SetTTL(key, val string, ttl int64) (resp *clientv3.PutResponse, err error) {
	cli, err := requireClient()
	if err != nil {
		return nil, err
	}
	return cli.SetTTL(key, val, ttl)
}

func SetTTLWithPrevKV(key, val string, ttl int64) (resp *clientv3.PutResponse, err error) {
	cli, err := requireClient()
	if err != nil {
		return nil, err
	}
	return cli.SetTTLWithPrevKV(key, val, ttl)
}

func Get(key string, opts ...clientv3.OpOption) (rsp *clientv3.GetResponse, err error) {
	cli, err := requireClient()
	if err != nil {
		return nil, err
	}
	return cli.Get(key, opts...)
}

func GetString(key string) string {
	cli, err := requireClient()
	if err != nil {
		return ""
	}
	return cli.GetString(key)
}

func GetArray(key string) (data []string) {
	cli, err := requireClient()
	if err != nil {
		return nil
	}
	return cli.GetArray(key)
}

func GetMap(key string) (data map[string]string) {
	cli, err := requireClient()
	if err != nil {
		return nil
	}
	return cli.GetMap(key)
}

func Del(key string) (resp *clientv3.DeleteResponse, err error) {
	cli, err := requireClient()
	if err != nil {
		return nil, err
	}
	return cli.Del(key)
}

func DelWithPrefix(key string) (resp *clientv3.DeleteResponse, err error) {
	cli, err := requireClient()
	if err != nil {
		return nil, err
	}
	return cli.DelWithPrefix(key)
}

func RegisterService(key, val string) (err error) {
	cli, err := requireClient()
	if err != nil {
		return err
	}
	return cli.RegisterService(key, val)
}

func Watch(key string) <-chan []string {
	cli, err := requireClient()
	if err != nil {
		ch := make(chan []string)
		close(ch)
		return ch
	}
	return cli.Watch(key)
}
