package util

import (
	"github.com/patrickmn/go-cache"
	"time"
)

func NewGoCache() *cache.Cache {
	return cache.New(time.Minute, 10*time.Minute)
}
