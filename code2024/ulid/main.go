/*

	ULID是一种基于UUID的通用唯一标识符，具有良好的排序性和可读性

*/

package main

import (
	"fmt"
	ulid "github.com/oklog/ulid/v2"
	"math/rand"
	"time"
)

func main() {
	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	ms := ulid.Timestamp(time.Now())
	for i := 0; i < 1000; i++ {
		obj, _ := ulid.New(ms, entropy)
		fmt.Println(obj.String())
	}

	for i := 0; i < 1000; i++ {
		fmt.Println(ulid.Make())
	}

}
