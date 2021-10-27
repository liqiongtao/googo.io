package goo_kv

type adapter interface {
	// 添加配置
	Set(key, val string, ttl int64) (err error)
	// 获取配置
	Get(key string) string
	// 获取配置
	GetMap(key string) (data map[string]string)
	// 删除配置
	Del(key string) (err error)
	// 监控配置
	Watch(key string)
	// 续约配置
	TTL(key string, ttl int64) (err error)
}
