package main

import (
	"chapter1/server/cache"
	"chapter1/server/http"
)

func main() {
	c := cache.New("inmemory") // 创建缓存
	http.New(c).Listen()
}
