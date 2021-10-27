package goo_kv

import (
	"context"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

type etcdAdapter struct {
	cli *clientv3.Client
	ctx context.Context
}

func NewEtcdAdapter(endpoints []string) *etcdAdapter {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err.Error())
	}
	return &etcdAdapter{
		cli: cli,
		ctx: context.Background(),
	}
}

func (ea etcdAdapter) Client() *clientv3.Client {
	return ea.cli
}

func (ea etcdAdapter) Set(key, val string, ttl int64) (err error) {
	var lease *clientv3.LeaseGrantResponse

	lease, err = ea.cli.Grant(ea.ctx, ttl)
	if err != nil {
		goo_log.Error(err.Error())
		return
	}

	_, err = ea.cli.Put(ea.ctx, key, val, clientv3.WithLease(lease.ID))
	if err != nil {
		goo_log.Error(err.Error())
	}
	return
}

func (ea etcdAdapter) Get(key string) string {
	rst, err := ea.cli.Get(ea.ctx, key, clientv3.WithPrefix())
	if err != nil {
		goo_log.Error(err.Error())
		return ""
	}
	if l := len(rst.Kvs); l == 0 {
		return ""
	}
	return string(rst.Kvs[0].Value)
}

func (ea etcdAdapter) GetMap(key string) (data map[string]string) {
	data = map[string]string{}

	rst, err := ea.cli.Get(ea.ctx, key, clientv3.WithPrefix())
	if err != nil {
		goo_log.Error(err.Error())
		return
	}
	if l := len(rst.Kvs); l == 0 {
		return
	}

	for _, i := range rst.Kvs {
		data[string(i.Key)] = string(i.Value)
	}
	return
}

func (ea etcdAdapter) Del(key string) (err error) {
	_, err = ea.cli.Delete(ea.ctx, key)
	if err != nil {
		goo_log.Error(err.Error())
	}
	return
}

func (ea etcdAdapter) Watch(key string) {
	//ch := client().Watch(ctx, key, clientv3.WithPrefix())
	//for {
	//	select {
	//	case rst := <-ch:
	//		for _, env := range rst.Events {
	//			fmt.Println(env.Type, string(env.Kv.Key), string(env.Kv.Value))
	//			//	fmt.Println("watch", key, env.Type, string(env.PrevKv.Key), string(env.PrevKv.Value))
	//		}
	//	case <-goo.CancelContext().Done():
	//		return
	//	}
	//}
}

func (ea etcdAdapter) TTL(key string, ttl int64) (err error) {
	lease := clientv3.NewLease(ea.cli)

	var rsp *clientv3.LeaseGrantResponse
	rsp, err = lease.Grant(ea.ctx, ttl)
	if err != nil {
		goo_log.Error(err.Error())
		return
	}

	if _, err := lease.KeepAliveOnce(ea.ctx, rsp.ID); err != nil {
		goo_log.Error(err.Error())
	}
	return
}
