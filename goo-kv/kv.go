package goo_kv

var __adapter adapter

func SetAdapter(adapter_ adapter) {
	__adapter = adapter_
}

// 添加配置
func Set(key, val string, ttl int64) (err error) {
	return __adapter.Set(key, val, ttl)
}

// 获取配置
func Get(key string) string {
	return __adapter.Get(key)
}

// 获取配置
func GetMap(key string) (data map[string]string) {
	return __adapter.GetMap(key)
}

// 删除配置
func Del(key string) (err error) {
	return __adapter.Del(key)
}

// 监控配置
func Watch(key string) {
	__adapter.Watch(key)
}

// 续约配置
func TTL(key string, ttl int64) (err error) {
	return __adapter.TTL(key, ttl)
}
