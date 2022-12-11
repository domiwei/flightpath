package localcache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type Service interface {
	SetDefault(key string, value interface{})
	Get(key string) (interface{}, bool)
	Delete(key string)
}

func NewLocalCache() Service {
	return cache.New(5*time.Minute, 10*time.Minute)
}
