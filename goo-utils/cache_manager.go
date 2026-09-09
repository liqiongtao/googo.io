package goo_utils

import (
	"github.com/liqiongtao/googo.io/goocontext"
	"golang.org/x/sync/singleflight"
	"sync"
	"time"
)

type CacheManager struct {
	duration time.Duration      // 缓存有效期
	cache    sync.Map           // 缓存数据
	sfGroup  singleflight.Group // 用于控制并发重复调用
}

type cacheItem struct {
	value  interface{} // 缓存值
	expire int64       // 过期时间戳（Unix 时间戳，单位：秒）
}

func NewCacheManager() *CacheManager {
	cm := &CacheManager{
		duration: time.Second * 300,
		cache:    sync.Map{},
		sfGroup:  singleflight.Group{},
	}
	// 启动后台清理 goroutine
	go cm.cleanupExpiredItems()
	return cm
}

func (cm *CacheManager) WithDuration(duration time.Duration) *CacheManager {
	cm.duration = duration
	return cm
}

func (cm *CacheManager) Get(key string, queryFunc func() (interface{}, error)) (interface{}, error) {
	// 1. 检查缓存
	if v, ok := cm.cache.Load(key); ok {
		item := v.(cacheItem)

		// 检查是否过期
		if item.expire > time.Now().Unix() {
			return item.value, nil
		}

		// 如果过期，删除缓存项
		cm.cache.Delete(key)
	}

	// 2. 使用 singleflight 确保相同 key 的并发请求只执行一次
	result, err, _ := cm.sfGroup.Do(key, func() (interface{}, error) {
		// 3. 再次检查缓存（可能其他 goroutine 已经查询并缓存了数据）
		if v, ok := cm.cache.Load(key); ok {
			item := v.(cacheItem)

			// 检查是否过期
			if item.expire > time.Now().Unix() {
				return item.value, nil
			}

			// 如果过期，删除缓存项
			cm.cache.Delete(key)
		}

		// 4. 缓存不存在，执行查询函数
		data, err := queryFunc()
		if err != nil {
			return nil, err
		}

		// 5. 缓存数据
		cm.cache.Store(key, cacheItem{
			value:  data,
			expire: time.Now().Unix() + int64(cm.duration.Seconds()),
		})

		return data, nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (cm *CacheManager) Set(key string, value interface{}) {
	cm.cache.Store(key, cacheItem{
		value:  value,
		expire: time.Now().Unix() + int64(cm.duration.Seconds()),
	})
}

func (cm *CacheManager) Delete(key string) {
	cm.cache.Delete(key)
}

func (cm *CacheManager) cleanupExpiredItems() {
	ticker := time.NewTicker(1 * time.Minute) // 每分钟清理一次
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			now := time.Now().Unix()
			cm.cache.Range(func(key, value interface{}) bool {
				item := value.(cacheItem)
				if item.expire > 0 && item.expire < now {
					cm.cache.Delete(key)
				}
				return true
			})
		case <-goocontext.Root().Done():
			return
		}
	}
}
