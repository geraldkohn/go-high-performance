package concurrentmap

import (
	"hash/crc32"
	"sync"
)

var default_shard_count = 64

// 分成default_shard_count个分片的map
type concurrentMap []*concurrentMapShard

// 通过RWMutex线程安全保护的分片, 包含一个map
type concurrentMapShard struct {
	sync.RWMutex
	items map[string]interface{}
}

// 创建并发Map
func New() *concurrentMap {
	m := make(concurrentMap, default_shard_count)
	for i := 0; i < len(m); i++ {
		m[i] = &concurrentMapShard{
			items: make(map[string]interface{}),
		}
	}
	return &m
}

// 根据key计算落在哪个分片上.
func (m *concurrentMap) getShard(key string) *concurrentMapShard {
	hashKey := hash(key)
	index := hashKey%uint32(default_shard_count)
	return (*m)[index]
}

func (m *concurrentMap) Get(key string) (interface{}, bool) {
	shard := m.getShard(key)
	shard.RLock()
	v, ok := shard.items[key]
	shard.RUnlock()
	return v, ok
}

func (m *concurrentMap) Set(key string, value interface{}) {
	shard := m.getShard(key)
	shard.Lock()
	shard.items[key] = value
	shard.Unlock()
}

// hash 算法, 映射到uint32上
func hash(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}
