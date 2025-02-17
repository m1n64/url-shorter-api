package utils

import "github.com/bradfitz/gomemcache/memcache"

func CreateMemcachedConn(host string, port string) *memcache.Client {
	return memcache.New(host + ":" + port)
}
