package goo_etcd

import (
	"context"
	"crypto/tls"
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"time"

	goo_log "github.com/liqiongtao/googo.io/goo-log"
	"github.com/liqiongtao/googo.io/goocontext"
	"go.etcd.io/etcd/client/pkg/v3/transport"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"go.uber.org/zap"
)

type Client struct {
	*clientv3.Client

	ctx  context.Context
	conf Config
}

func New(conf Config) (cli *Client, err error) {
	cli = &Client{ctx: goocontext.Root(), conf: conf}

	cfg := clientv3.Config{
		Endpoints:   conf.Endpoints,
		DialTimeout: 5 * time.Second,
		Logger:      zap.NewNop(),
	}

	if conf.Username != "" {
		cfg.Username = conf.Username
	}
	if conf.Password != "" {
		cfg.Password = conf.Password
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
			goo_log.WithTag("goo-etcd").WithField("config", conf).Error(err.Error())
			return nil, err
		}
		cfg.TLS = clientConfig
	}

	cli.Client, err = clientv3.New(cfg)
	if err != nil {
		goo_log.WithTag("goo-etcd").WithField("config", conf).Error(err.Error())
		return nil, err
	}

	return cli, nil
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

	opts = append(opts, clientv3.WithLease(lease.ID))
	resp, err = cli.Client.Put(cli.ctx, key, val, opts...)
	if err != nil {
		_, _ = cli.Client.Revoke(cli.ctx, lease.ID)
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
	if l := len(resp.Kvs); l == 0 {
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

// register service and keepalive
func (cli *Client) RegisterService(serviceName, addr string) (err error) {
	for {
		if err = cli.ctx.Err(); err != nil {
			return err
		}

		err = cli.registerServiceOnce(serviceName, addr)
		if err == nil {
			return nil
		}

		select {
		case <-cli.ctx.Done():
			return cli.ctx.Err()
		case <-time.After(3 * time.Second):
		}
	}
}

func (cli *Client) registerServiceOnce(serviceName, addr string) (err error) {
	if cli.Client == nil {
		return fmt.Errorf("etcd client is nil")
	}

	var (
		ttl   int64 = 5
		em    endpoints.Manager
		lease *clientv3.LeaseGrantResponse
		ch    <-chan *clientv3.LeaseKeepAliveResponse
	)

	lease, err = cli.Client.Grant(cli.ctx, ttl)
	if err != nil {
		goo_log.WithTag("goo-etcd").WithField("serviceName", serviceName).WithField("addr", addr).Error(err)
		return
	}

	em, err = endpoints.NewManager(cli.Client, serviceName)
	if err != nil {
		goo_log.WithTag("goo-etcd").WithField("serviceName", serviceName).WithField("addr", addr).Error(err)
		return
	}

	serviceKey := serviceName + "/" + strconv.Itoa(int(lease.ID))
	err = em.AddEndpoint(cli.ctx, serviceKey, endpoints.Endpoint{Addr: addr}, clientv3.WithLease(lease.ID))
	if err != nil {
		goo_log.WithTag("goo-etcd").WithField("serviceName", serviceName).WithField("addr", addr).Error(err)
		return
	}

	ch, err = cli.Client.KeepAlive(cli.ctx, lease.ID)
	if err != nil {
		goo_log.WithTag("goo-etcd").WithField("serviceName", serviceName).WithField("addr", addr).Error(err)
		return
	}

	goo_log.WithTag("goo-etcd").WithField("serviceName", serviceName).WithField("addr", addr).Debug("服务注册成功")

	go func() {
		for {
			select {
			case <-cli.ctx.Done():
				goo_log.WithTag("goo-etcd").WithField("serviceName", serviceName).WithField("addr", addr).Warn("服务退出,收回注册信息")
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				_, _ = cli.Client.Revoke(ctx, lease.ID)
				cancel()
				return

			case rsp, ok := <-ch:
				if !ok || rsp == nil {
					goo_log.WithTag("goo-etcd").WithField("serviceName", serviceName).WithField("addr", addr).Error("服务注册续租失效")
					go func() {
						_ = cli.RegisterService(serviceName, addr)
					}()
					return
				}
			}
		}
	}()

	return
}

// watch the key
func (cli *Client) Watch(key string) <-chan []string {
	var (
		mu   sync.Mutex
		ch   = make(chan []string, runtime.NumCPU()*2)
		data = cli.GetMap(key)
	)

	ch <- cli.map2array(data)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				goo_log.WithTag("goo-etcd").WithField("key", key).Error(r)
			}
			close(ch)
		}()

		wc := cli.Client.Watch(cli.ctx, key, clientv3.WithPrefix())
		for {
			select {
			case <-cli.ctx.Done():
				return

			case w, ok := <-wc:
				if !ok {
					return
				}
				if err := w.Err(); err != nil {
					goo_log.WithTag("goo-etcd").WithField("key", key).Error(err)
					continue
				}

				mu.Lock()
				for _, ev := range w.Events {
					k := string(ev.Kv.Key)
					v := string(ev.Kv.Value)

					switch ev.Type {
					case clientv3.EventTypePut:
						data[k] = v

					case clientv3.EventTypeDelete:
						delete(data, k)
					}
				}
				arr := cli.map2array(data)
				mu.Unlock()

				select {
				case ch <- arr:
				case <-cli.ctx.Done():
					return
				}
			}
		}
	}()

	return ch
}

func (cli *Client) map2array(data map[string]string) []string {
	var arrData []string
	for _, v := range data {
		arrData = append(arrData, v)
	}
	return arrData
}
