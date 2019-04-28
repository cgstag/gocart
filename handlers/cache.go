package handlers

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var Cache *cache.Cache

func SetupCache(defaultExpiration int, defaultPurge int) error {
	Cache = cache.New(time.Duration(defaultExpiration)*time.Minute, time.Duration(defaultPurge)*time.Minute)
	return nil
}
