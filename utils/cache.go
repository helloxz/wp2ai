package utils

import "github.com/coocood/freecache"

var Cache *freecache.Cache

func InitCache() {
	// 设置 5MB 的缓存大小
	cacheSize := 5 * 1024 * 1024
	Cache = freecache.NewCache(cacheSize)
}
