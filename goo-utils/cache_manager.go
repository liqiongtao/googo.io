package goo_utils

import (
	"golang.org/x/sync/singleflight"
	"sync"
)

type CacheManager struct {
	v       sync.Map           // 缓存
	sfGroup singleflight.Group // 用于控制并发重复调用
}

func NewCacheManager() *CacheManager {
	return &CacheManager{
		v:       sync.Map{},
		sfGroup: singleflight.Group{},
	}
}

func (cm *CacheManager) Get(key string, queryFunc func() (interface{}, error)) (interface{}, error) {
	// 1. 检查缓存
	if v, ok := cm.v.Load(key); ok {
		return v, nil
	}

	// 2. 使用 singleflight 确保相同 key 的并发请求只执行一次
	result, err, _ := cm.sfGroup.Do(key, func() (interface{}, error) {
		// 3. 再次检查缓存（可能其他 goroutine 已经查询并缓存了数据）
		if v, ok := cm.v.Load(key); ok {
			return v, nil
		}

		// 4. 缓存不存在，执行查询函数
		data, err := queryFunc()
		if err != nil {
			return nil, err
		}

		// 5. 缓存数据
		cm.v.Store(key, data)
		return data, nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (cm *CacheManager) Update(key string, value interface{}) {
	cm.v.Store(key, value)
}

func (cm *CacheManager) Delete(key string) {
	cm.v.Delete(key)
}
