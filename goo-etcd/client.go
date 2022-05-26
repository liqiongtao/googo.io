package goo_etcd

import (
	"context"
	"crypto/tls"
	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"go.etcd.io/etcd/client/pkg/v3/transport"
	clientv3 "go.etcd.io/etcd/client/v3"
	"runtime"
	"time"
)

type Client struct {
	*clientv3.Client
	ctx context.Context
}

func New(conf Config) (cli *Client, err error) {
	cli = &Client{ctx: context.TODO()}

	config := clientv3.Config{
		Endpoints:   conf.Endpoints,
		DialTimeout: 5 * time.Second,
	}

	if conf.Username != "" {
		config.Username = conf.Username
	}
	if conf.Password != "" {
		config.Password = conf.Password
	}

	if conf.TLS != nil {
		tlsInfo := &transport.TLSInfo{
			CertFile:      conf.TLS.CertFile,
			KeyFile:       conf.TLS.KeyFile,
			TrustedCAFile: conf.TLS.CAFile,
		}
		var clientConfig *tls.Config
		clientConfig, err = tlsInfo.ClientConfig()
		if err != nil {
			goo_log.WithTag("goo-etcd").Error(err.Error())
			return
		}
		config.TLS = clientConfig
	}

	cli.Client, err = clientv3.New(config)
	if err != nil {
		goo_log.WithTag("goo-etcd").Error(err.Error())
	}
	return
}

// set key-value
func (cli *Client) Set(key, val string, opts ...clientv3.OpOption) (resp *clientv3.PutResponse, err error) {
	resp, err = cli.Client.Put(cli.ctx, key, val, opts...)
	if err != nil {
		goo_log.WithTag("goo-etcd").WithField("key", key).WithField("val", val).Error(err)
	}
	return
}

// set key-value and return previous key-value
func (cli *Client) SetWithPrevKV(key, val string) (resp *clientv3.PutResponse, err error) {
	return cli.Set(key, val, clientv3.WithPrevKV())
}

// set key-value-ttl
func (cli *Client) SetTTL(key, val string, ttl int64, opts ...clientv3.OpOption) (resp *clientv3.PutResponse, err error) {
	if ttl == 0 {
		return cli.Set(key, val, opts...)
	}

	var lease *clientv3.LeaseGrantResponse

	lease, err = cli.Client.Grant(cli.ctx, ttl)
	if err != nil {
		goo_log.WithTag("goo-etcd").WithField("key", key).WithField("val", val).WithField("ttl", ttl).Error(err)
		return
	}

	_, err = cli.Client.Put(cli.ctx, key, val, clientv3.WithLease(lease.ID))
	if err != nil {
		goo_log.WithTag("goo-etcd").WithField("key", key).WithField("val", val).WithField("ttl", ttl).Error(err)
		return
	}

	return
}

// set key-value-ttl and return previous key-value
func (cli *Client) SetTTLWithPrevKV(key, val string, ttl int64) (resp *clientv3.PutResponse, err error) {
	return cli.SetTTL(key, val, ttl, clientv3.WithPrevKV())
}

// get value by key
func (cli *Client) Get(key string, opts ...clientv3.OpOption) (resp *clientv3.GetResponse, err error) {
	resp, err = cli.Client.Get(cli.ctx, key, opts...)
	if err != nil {
		goo_log.WithTag("goo-etcd").WithField("key", key).Error(err)
	}
	return
}

// get string value by prefix key
func (cli *Client) GetString(key string) string {
	resp, err := cli.Get(key, clientv3.WithPrefix())
	if err != nil {
		return ""
	}
	return string(resp.Kvs[0].Value)
}

// get array value by prefix key
func (cli *Client) GetArray(key string) (data []string) {
	data = []string{}

	resp, err := cli.Get(key, clientv3.WithPrefix())
	if err != nil {
		goo_log.WithTag("goo-etcd").WithField("key", key).Error(err)
		return
	}

	for _, i := range resp.Kvs {
		data = append(data, string(i.Value))
	}

	return
}

// get map value by prefix key
func (cli *Client) GetMap(key string) (data map[string]string) {
	data = map[string]string{}

	resp, err := cli.Get(key, clientv3.WithPrefix())
	if err != nil {
		goo_log.WithTag("goo-etcd").WithField("key", key).Error(err)
		return
	}

	for _, i := range resp.Kvs {
		key := string(i.Key)
		data[key] = string(i.Value)
	}

	return
}

// del key and return previous key-value
func (cli *Client) Del(key string, opts ...clientv3.OpOption) (resp *clientv3.DeleteResponse, err error) {
	opts = append(opts, clientv3.WithPrevKV())
	resp, err = cli.Client.Delete(cli.ctx, key, opts...)
	if err != nil {
		goo_log.WithTag("goo-etcd").WithField("key", key).Error(err)
	}
	return
}

// del prefix key and return previous key-value
func (cli *Client) DelWithPrefix(key string) (resp *clientv3.DeleteResponse, err error) {
	return cli.Del(key, clientv3.WithPrefix())
}

// set key and keep alive
func (cli *Client) RegisterService(key, val string) (err error) {
	var (
		ttl   int64 = 1
		lease *clientv3.LeaseGrantResponse
		ch    <-chan *clientv3.LeaseKeepAliveResponse
	)

	lease, err = cli.Client.Grant(cli.ctx, ttl)
	if err != nil {
		goo_log.WithTag("goo-etcd").WithField("key", key).WithField("val", val).Error(err)
		return
	}

	leaseID := lease.ID

	_, err = cli.Client.Put(cli.ctx, key, val, clientv3.WithLease(leaseID))
	if err != nil {
		goo_log.WithTag("goo-etcd").WithField("key", key).WithField("val", val).Error(err)
		return
	}

	ch, err = cli.Client.KeepAlive(cli.ctx, leaseID)
	if err != nil {
		goo_log.WithTag("goo-etcd").WithField("key", key).WithField("val", val).Error(err)
		return
	}

	go func() {
		for {
			select {
			case <-ch:
				return

			case <-cli.ctx.Done():
				return
			}
		}
	}()

	return
}

// watch the key
func (cli *Client) Watch(key string) <-chan []string {
	ch := make(chan []string, runtime.NumCPU()*2)

	go func() {
		defer func() {
			if err := recover(); err != nil {
				goo_log.WithTag("goo-etcd").WithField("key", key).Error(err)
			}
		}()

		ch <- cli.GetArray(key)
	}()

	go func() {
		defer func() {
			if err := recover(); err != nil {
				goo_log.WithTag("goo-etcd").WithField("key", key).Error(err)
			}
		}()

		wc := cli.Client.Watch(cli.ctx, key, clientv3.WithPrefix())
		for range wc {
			ch <- cli.GetArray(key)
		}
	}()

	return ch
}
