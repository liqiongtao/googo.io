package goo_etcd

import (
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

var __client *Client

func Init(conf Config) {
	var err error
	__client, err = New(conf)

	goo_utils.AsyncFunc(func() {
		for {
			select {
			case <-__client.ctx.Done():
				__client.Close()
				time.Sleep(time.Second)
				return

			default:
				if err != nil {
					time.Sleep(5 * time.Second)
					__client, err = New(conf)
				}
			}
		}
	})
}

func Default() *Client {
	return __client
}

func Set(key, val string) (resp *clientv3.PutResponse, err error) {
	return __client.Set(key, val)
}

func SetWithPrevKV(key, val string) (resp *clientv3.PutResponse, err error) {
	return __client.SetWithPrevKV(key, val)
}

func SetTTL(key, val string, ttl int64) (resp *clientv3.PutResponse, err error) {
	return __client.SetTTL(key, val, ttl)
}

func SetTTLWithPrevKV(key, val string, ttl int64) (resp *clientv3.PutResponse, err error) {
	return __client.SetTTLWithPrevKV(key, val, ttl)
}

func Get(key string, opts ...clientv3.OpOption) (rsp *clientv3.GetResponse, err error) {
	return __client.Get(key, opts...)
}

func GetString(key string) string {
	return __client.GetString(key)
}

func GetArray(key string) (data []string) {
	return __client.GetArray(key)
}

func GetMap(key string) (data map[string]string) {
	return __client.GetMap(key)
}

func Del(key string) (resp *clientv3.DeleteResponse, err error) {
	return __client.Del(key)
}

func DelWithPrefix(key string) (resp *clientv3.DeleteResponse, err error) {
	return __client.DelWithPrefix(key)
}

func RegisterService(key, val string) (err error) {
	return __client.RegisterService(key, val)
}

func Watch(key string) <-chan []string {
	return __client.Watch(key)
}
