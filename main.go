package main

import (
	"fmt"

	"git.zituo.net/liman/mange/core/manredis"
)

func main() {
	redis_host := "127.0.0.1"
	redis_port := 6379
	redis_password := ""

	r := &manredis.RedisGo{
		RedisIP:       redis_host,
		RedisPort:     redis_port,
		RedisPassword: redis_password,
		ConnType:      1,
	}

	c := r.RConn()

	fmt.Println(c)

}
