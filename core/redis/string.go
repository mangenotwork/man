package manredis

import (
	"fmt"

	_ "github.com/garyburd/redigo/redis"
)

type StringCMD struct{}

//	SET
//	SET key value [EX seconds] [PX milliseconds] [NX|XX]
//	将字符串值 value 关联到 key 。
//	如果 key 已经持有其他值， SET 就覆写旧值，无视类型。
//	对于某个原本带有生存时间（TTL）的键来说， 当 SET 命令成功在这个键上执行时， 这个键原有的 TTL 将被清除。
// func (r *StringCMD) StrSET(key string, value interface{}) (err error) {
// 	fmt.Println("[Execute redis command]: ", "SET", key, value)
// 	conn := Getconn()

// 	_, err = conn.Do("SET", key, value)
// 	return
// }
