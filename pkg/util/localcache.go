package util

import (
	"github.com/dgraph-io/ristretto"
	"math/rand"
)

//本地缓存

type RistrettoCache struct {
	cache *ristretto.Cache
}

func NewRistrettoCache() (*RistrettoCache, error) {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters:            1e7,     // 跟踪的键数量 (10M)
		MaxCost:                1 << 30, // 最大缓存大小 (1GB)
		BufferItems:            64,      // 优化性能
		TtlTickerDurationInSec: 6*60*2 + int64(rand.Intn(30*2)),
	})

	if err != nil {
		return nil, err
	}

	return &RistrettoCache{cache: cache}, nil
}

func (rc *RistrettoCache) Set(key, value interface{}, cost int64) {
	rc.cache.Set(key, value, cost)
}

func (rc *RistrettoCache) Get(key interface{}) (interface{}, bool) {
	return rc.cache.Get(key)
}

func (rc *RistrettoCache) Del(key interface{}) {
	rc.cache.Del(key)
}
