package gooredis

import (
	"context"
	"time"

	goocontext "github.com/liqiongtao/googo.io/goo-context"
	goolog "github.com/liqiongtao/googo.io/goo-log"
	"github.com/redis/go-redis/v9"
)

type Client struct {
	Config
	*redis.Client
}

func New(conf Config) (cli *Client, err error) {
	cli = &Client{Config: conf}

	opts := &redis.Options{
		Addr:     conf.Addr,
		Password: conf.Password,
		DB:       conf.DB,

		// 连接池
		PoolSize:        20,
		MinIdleConns:    5,
		PoolTimeout:     30 * time.Second,
		ConnMaxIdleTime: 5 * time.Minute,

		// 超时
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	if o := conf.Options; o != nil {
		if o.Addr != "" {
			opts.Addr = o.Addr
		}
		if o.Password != "" {
			opts.Password = o.Password
		}
		if o.DB != 0 {
			opts.DB = o.DB
		}
		if o.PoolSize > 0 {
			opts.PoolSize = o.PoolSize
		}
		if o.MinIdleConns > 0 {
			opts.MinIdleConns = o.MinIdleConns
		}
		if o.PoolTimeout > 0 {
			opts.PoolTimeout = o.PoolTimeout
		}
		if o.ConnMaxIdleTime > 0 {
			opts.ConnMaxIdleTime = o.ConnMaxIdleTime
		}
		if o.DialTimeout > 0 {
			opts.DialTimeout = o.DialTimeout
		}
		if o.ReadTimeout > 0 {
			opts.ReadTimeout = o.ReadTimeout
		}
		if o.WriteTimeout > 0 {
			opts.WriteTimeout = o.WriteTimeout
		}
		if o.Username != "" {
			opts.Username = o.Username
		}
		if o.TLSConfig != nil {
			opts.TLSConfig = o.TLSConfig
		}
		if o.Dialer != nil {
			opts.Dialer = o.Dialer
		}
		if o.OnConnect != nil {
			opts.OnConnect = o.OnConnect
		}
	}

	cli.Client = redis.NewClient(opts)

	if err = cli.Client.Ping(cli.cmdContext()).Err(); err != nil {
		goolog.WithTag("goo-redis").Error(err)
		_ = cli.Close()
		cli = nil
		return
	}

	return
}

// cmdContext 单次命令的 ctx。不用 Client 生命周期 Context：
// Exit 后 Root 已 cancel，但收尾写（删 key、重入队等）仍须成功。
// 阻塞订阅见 Subscribe（跟 goocontext.Root）。
func (c *Client) cmdContext() context.Context {
	return context.Background()
}

// ---- 兼容旧调用（无 context 参数）----

func (c *Client) Ping() *redis.StatusCmd {
	return c.Client.Ping(c.cmdContext())
}

func (c *Client) Get(key string) *redis.StringCmd {
	return c.Client.Get(c.cmdContext(), key)
}

func (c *Client) Set(key string, value any, expiration time.Duration) *redis.StatusCmd {
	return c.Client.Set(c.cmdContext(), key, value, expiration)
}

func (c *Client) SetNX(key string, value any, expiration time.Duration) *redis.BoolCmd {
	return c.Client.SetNX(c.cmdContext(), key, value, expiration)
}

func (c *Client) Del(keys ...string) *redis.IntCmd {
	return c.Client.Del(c.cmdContext(), keys...)
}

func (c *Client) Exists(keys ...string) *redis.IntCmd {
	return c.Client.Exists(c.cmdContext(), keys...)
}

func (c *Client) Expire(key string, expiration time.Duration) *redis.BoolCmd {
	return c.Client.Expire(c.cmdContext(), key, expiration)
}

func (c *Client) Eval(script string, keys []string, args ...any) *redis.Cmd {
	return c.Client.Eval(c.cmdContext(), script, keys, args...)
}

func (c *Client) Subscribe(channels ...string) *redis.PubSub {
	return c.Client.Subscribe(goocontext.Root(), channels...)
}

func (c *Client) HSet(key string, values ...any) *redis.IntCmd {
	return c.Client.HSet(c.cmdContext(), key, values...)
}

func (c *Client) HGet(key, field string) *redis.StringCmd {
	return c.Client.HGet(c.cmdContext(), key, field)
}

func (c *Client) HGetAll(key string) *redis.MapStringStringCmd {
	return c.Client.HGetAll(c.cmdContext(), key)
}

func (c *Client) HDel(key string, fields ...string) *redis.IntCmd {
	return c.Client.HDel(c.cmdContext(), key, fields...)
}

func (c *Client) HKeys(key string) *redis.StringSliceCmd {
	return c.Client.HKeys(c.cmdContext(), key)
}

func (c *Client) HMSet(key string, values ...any) *redis.BoolCmd {
	return c.Client.HMSet(c.cmdContext(), key, values...)
}

func (c *Client) HIncrBy(key, field string, incr int64) *redis.IntCmd {
	return c.Client.HIncrBy(c.cmdContext(), key, field, incr)
}

func (c *Client) ZAdd(key string, members ...redis.Z) *redis.IntCmd {
	return c.Client.ZAdd(c.cmdContext(), key, members...)
}

func (c *Client) ZRem(key string, members ...any) *redis.IntCmd {
	return c.Client.ZRem(c.cmdContext(), key, members...)
}

func (c *Client) ZCard(key string) *redis.IntCmd {
	return c.Client.ZCard(c.cmdContext(), key)
}

func (c *Client) ZRange(key string, start, stop int64) *redis.StringSliceCmd {
	return c.Client.ZRange(c.cmdContext(), key, start, stop)
}

func (c *Client) ZRangeWithScores(key string, start, stop int64) *redis.ZSliceCmd {
	return c.Client.ZRangeWithScores(c.cmdContext(), key, start, stop)
}

func (c *Client) Publish(channel string, message any) *redis.IntCmd {
	return c.Client.Publish(c.cmdContext(), channel, message)
}

func (c *Client) TxPipeline() redis.Pipeliner {
	return c.Client.TxPipeline()
}
