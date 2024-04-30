/*
	泛型
*/

package main

import (
	"golang.org/x/exp/constraints"
	"log"
)

func main() {

	test1 := []int{1, 1, 2, 2, 3, 3}
	test2 := []string{"1", "1", "2", "2", "3", "3"}

	log.Println(SliceRemoveDuplicate(test1))
	log.Println(SliceRemoveDuplicate(test2))

	log.Println(Min(1, 3))
	log.Println(Min(3.2, 3))

	log.Println(IsContain(test1, 1))
	log.Println(IsContain(test2, "11"))

	log.Println(IF(true, "1", "asdas"))
	log.Println(IF(false, 1, 1111))

	test11 := CopySlice(test1)
	log.Println(test11)
	test22 := CopySlice(test2)
	log.Println(test22)

}

type SliceType interface {
	int | uint | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64 | float32 | float64 | string | bool |
	*int | *uint | *int8 | *uint8 | *int16 | *uint16 | *int32 | *uint32 | *int64 | *uint64 | *float32 | *float64 | *string | *bool |
	chan int | chan uint | chan int8 | chan uint8 | chan int16 | chan uint16 | chan int32 | chan uint32 | chan int64 | chan uint64 | chan float32 | chan float64 | chan string | chan bool |
	chan *int | chan *uint | chan *int8 | chan *uint8 | chan *int16 | chan *uint16 | chan *int32 | chan *uint32 | chan *int64 | chan *uint64 | chan *float32 | chan *float64 | chan *string | chan *bool
}

// SliceRemoveDuplicate slice去重
func SliceRemoveDuplicate[T SliceType](arr []T) []T {
	result := make([]T, 0, len(arr))
	temp := make(map[T]struct{})
	for _, item := range arr {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

// SliceRemoveSpecificParam 从slice中移除特定的参数
func SliceRemoveSpecificParam[T SliceType](arr []T, specific T) []T {
	temp := make([]T, 0, len(arr))
	for _, value := range arr {
		if value != specific {
			temp = append(temp, value)
		}
	}
	return temp
}

// SliceIntersection slice的交集
func SliceIntersection[T SliceType](first []T, last []T) []T {
	result := make([]T, 0)
	temp := make(map[T]struct{})
	for _, value := range first {
		if _, ok := temp[value]; !ok {
			temp[value] = struct{}{}
		}
	}
	for _, val := range last {
		if _, ok := temp[val]; ok {
			result = append(result, val)
		}
	}
	return result
}

// SliceComplement slice的补集(subset必须是complete的子集)
func SliceComplement[T SliceType](complete []T, subset []T) []T {
	result := make([]T, 0)
	temp := make(map[T]struct{})
	for _, value := range subset {
		if _, ok := temp[value]; !ok {
			temp[value] = struct{}{}
		}
	}
	for _, val := range complete {
		if _, ok := temp[val]; !ok {
			temp[val] = struct{}{} // 去重
			result = append(result, val)
		}
	}
	return result
}

// SlicePaginate slice 分页(1开始)
func SlicePaginate[T SliceType](x []T, page int, size int) []T {
	skip := (page - 1) * size
	if skip > len(x) {
		skip = len(x)
	}
	end := skip + size
	if end > len(x) {
		end = len(x)
	}
	return x[skip:end]
}

// SliceIsExist 指定参数值是否存在data[slice]中 false:不存在
func SliceIsExist[T SliceType](data []T, param T) bool {
	temp := make(map[T]struct{})
	for _, value := range data {
		if _, ok := temp[value]; !ok {
			temp[value] = struct{}{}
		}
	}
	if _, ok := temp[param]; !ok {
		return false
	}
	return true
}

func isNaN[T constraints.Ordered](x T) bool {
	return x != x
}

func Min[T constraints.Ordered](a, b T) T {
	if a < b || isNaN(a) {
		return a
	}
	return b
}

func Max[T constraints.Ordered](a, b T) T {
	if a > b || isNaN(a) {
		return a
	}
	return b
}

func IsContain[T SliceType](items []T, item T) bool {
	for i := 0; i < len(items); i++ {
		if items[i] == item {
			return true
		}
	}
	return false
}

func IF[T any](condition bool, a, b T) T {
	if condition {
		return a
	}
	return b
}

func CopySlice[T SliceType](s []T) []T {
	return append(s[:0:0], s...)
}
