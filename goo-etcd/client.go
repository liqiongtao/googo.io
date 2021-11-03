package goo_etcd

import (
	"context"
	goo_context "github.com/liqiongtao/googo.io/goo-context"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	goo_utils "github.com/liqiongtao/googo.io/goo-utils"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Client struct {
	*clientv3.Client
	ctx context.Context
}

func (cli *Client) Set(key, val string, ttl int64) (err error) {
	var options []clientv3.OpOption

	if ttl > 0 {
		var lease *clientv3.LeaseGrantResponse

		lease, err = cli.Client.Grant(cli.ctx, ttl)
		if err != nil {
			goo_log.WithField("key", key).WithField("val", val).WithField("ttl", ttl).Error(err.Error())
			return
		}

		options = append(options, clientv3.WithLease(lease.ID))
	}

	_, err = cli.Client.Put(cli.ctx, key, val, options...)
	if err != nil {
		goo_log.WithField("key", key).WithField("val", val).WithField("ttl", ttl).Error(err.Error())
	}
	return
}

func (cli *Client) SetWithKeepAlive(key, val string, ttl int64) (err error) {
	if ttl == 0 {
		return
	}

	var (
		leaseGrantRsp     *clientv3.LeaseGrantResponse
		leaseKeepAliveRsp <-chan *clientv3.LeaseKeepAliveResponse
	)

	leaseGrantRsp, err = cli.Client.Grant(cli.ctx, ttl)
	if err != nil {
		goo_log.WithField("key", key).WithField("val", val).WithField("ttl", ttl).Error(err.Error())
		return
	}

	leaseId := leaseGrantRsp.ID

	_, err = cli.Client.Put(cli.ctx, key, val, clientv3.WithLease(leaseId))
	if err != nil {
		goo_log.WithField("key", key).WithField("val", val).WithField("ttl", ttl).Error(err.Error())
		return
	}

	goo_log.WithField("key", key).WithField("val", val).WithField("ttl", ttl).Debug("注册成功")

	leaseKeepAliveRsp, err = cli.Client.KeepAlive(cli.ctx, leaseId)
	if err != nil {
		goo_log.WithField("key", key).WithField("val", val).WithField("ttl", ttl).Error(err.Error())
		return
	}

	goo_utils.AsyncFunc(func() {
		for {
			select {
			case <-goo_context.Cancel().Done():
				return

			case rst := <-leaseKeepAliveRsp:
				if rst == nil {
					continue
				}
				goo_log.WithField("result", rst).Debug("keep-alive")
			}
		}
	})

	return
}

func (cli *Client) Get(key string) string {
	rsp, err := cli.Client.Get(cli.ctx, key)
	if err != nil {
		goo_log.WithField("key", key).Error(err.Error())
		return ""
	}
	if rsp.Count == 0 {
		return ""
	}
	return string(rsp.Kvs[0].Value)
}

func (cli *Client) GetMap(key string) (data map[string]string) {
	data = map[string]string{}
	rsp, err := cli.Client.Get(cli.ctx, key, clientv3.WithPrefix())
	if err != nil {
		goo_log.WithField("key", key).Error(err.Error())
		return
	}
	if rsp.Count == 0 {
		return
	}
	for _, kv := range rsp.Kvs {
		data[string(kv.Key)] = string(kv.Value)
	}
	return
}

func (cli *Client) GetArray(key string) (data []string) {
	data = []string{}
	rsp, err := cli.Client.Get(cli.ctx, key, clientv3.WithPrefix())
	if err != nil {
		goo_log.WithField("key", key).Error(err.Error())
		return
	}
	if rsp.Count == 0 {
		return
	}
	for _, kv := range rsp.Kvs {
		data = append(data, string(kv.Value))
	}
	return
}

func (cli *Client) Watch(key string, fn func([]string)) {
	ch := make(chan []string)

	goo_utils.AsyncFunc(func() {
		for {
			select {
			case <-goo_context.Cancel().Done():
				return
			case arr := <-ch:
				fn(arr)
			}
		}
	})

	ch <- cli.GetArray(key)

	goo_utils.AsyncFunc(func() {
		wch := cli.Client.Watch(cli.ctx, key, clientv3.WithPrefix())
		for rst := range wch {
			if rst.Canceled {
				close(ch)
				return
			}

			var arr []string
			for _, evt := range rst.Events {
				switch evt.Type {
				case mvccpb.PUT:
					arr = append(arr, string(evt.Kv.Value))
				}
			}

			ch <- arr
		}
	})
}

func (cli *Client) Del(key string) (err error) {
	_, err = cli.Client.Delete(cli.ctx, key)
	if err != nil {
		goo_log.WithField("key", key).Error(err.Error())
	}
	return
}
