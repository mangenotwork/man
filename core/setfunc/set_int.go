//	[]int  []int64  的setfunc

package setfunc

import (
	"sync"
)

// []int  去重
func SetInt(data *[]int) {
	l := new(sync.RWMutex)
	tempMap := make(map[int]bool, 0)
	for _, v := range *data {
		l.Lock()
		tempMap[v] = true
		l.Unlock()
	}
	*data = nil
	for k, _ := range tempMap {
		l.RLock()
		*data = append(*data, k)
		l.RUnlock()
	}
}

// []int64  去重
func SetInt64(data *[]int64) {
	l := new(sync.RWMutex)
	tempMap := make(map[int64]bool, 0)
	for _, v := range *data {
		l.Lock()
		tempMap[v] = true
		l.Unlock()
	}
	*data = nil
	for k, _ := range tempMap {
		l.RLock()
		*data = append(*data, k)
		l.RUnlock()
	}
}

// []int 并集  Merge
// 返回的是无序切片
func SetMergeInt(aSet, bSet []int) []int {
	l := new(sync.RWMutex)
	tempMap := make(map[int]bool, 0)
	for _, v := range aSet {
		l.Lock()
		tempMap[v] = true
		l.Unlock()
	}
	for _, v := range bSet {
		l.Lock()
		tempMap[v] = true
		l.Unlock()
	}
	U := make([]int, 0)
	for k, _ := range tempMap {
		l.RLock()
		U = append(U, k)
		l.RUnlock()
	}
	return U
}

// []int64 并集  Merge
// 返回的是无序切片
func SetMergeInt(aSet, bSet []int64) []int64 {
	l := new(sync.RWMutex)
	tempMap := make(map[int64]bool, 0)
	for _, v := range aSet {
		l.Lock()
		tempMap[v] = true
		l.Unlock()
	}
	for _, v := range bSet {
		l.Lock()
		tempMap[v] = true
		l.Unlock()
	}
	U := make([]int64, 0)
	for k, _ := range tempMap {
		l.RLock()
		U = append(U, k)
		l.RUnlock()
	}
	return U
}

// []int 交集 Unite
func SetUniteInt(aSet, bSet []int) []int {
	l := new(sync.RWMutex)
	tempMap := make(map[int]bool, 0)
	for _, v := range aSet {
		l.Lock()
		tempMap[v] = true
		l.Unlock()
	}
	U := make([]int, 0)
	for _, v := range bSet {
		l.RLock()
		if tempMap[v] {
			U = append(U, v)
		}
		l.RUnlock()
	}
	return U
}

// []int64 交集 Unite
func SetUniteInt(aSet, bSet []int64) []int64 {
	l := new(sync.RWMutex)
	tempMap := make(map[int64]bool, 0)
	for _, v := range aSet {
		l.Lock()
		tempMap[v] = true
		l.Unlock()
	}
	U := make([]int64, 0)
	for _, v := range bSet {
		l.RLock()
		if tempMap[v] {
			U = append(U, v)
		}
		l.RUnlock()
	}
	return U
}

// []int 差集 difference
func SetDifference(aSet, bSet []int) []int {
	l := new(sync.RWMutex)
	tempMap := make(map[int]bool, 0)
	for _, v := range aSet {
		l.Lock()
		tempMap[v] = true
		l.Unlock()
	}
	U := make([]int, 0)
	for _, v := range bSet {
		l.RLock()
		if !tempMap[v] {
			U = append(U, v)
		}
		l.RUnlock()
	}
	SetInt(&U)
	return U
}

// []int64 差集 difference
func SetDifference(aSet, bSet []int64) []int64 {
	l := new(sync.RWMutex)
	tempMap := make(map[int64]bool, 0)
	for _, v := range aSet {
		l.Lock()
		tempMap[v] = true
		l.Unlock()
	}
	U := make([]int64, 0)
	for _, v := range bSet {
		l.RLock()
		if !tempMap[v] {
			U = append(U, v)
		}
		l.RUnlock()
	}
	SetInt(&U)
	return U
}

// int 元素是否是集合的子集
func SubSetInt(val int, data []int) bool {
	for _, v := range data {
		if val == v {
			return true
		}
	}
	return false
}

// int64 元素是否是集合的子集
func SubSetInt(val int64, data []int64) bool {
	for _, v := range data {
		if val == v {
			return true
		}
	}
	return false
}

// []int 	b切片 是否包含 a切片的所有元素
// a∈b
func AllSubSetInt(aSet, bSet []int) bool {
	l := new(sync.RWMutex)
	tempMap := make(map[int]bool, 0)
	for _, v := range bSet {
		l.Lock()
		tempMap[v] = true
		l.Unlock()
	}
	for _, v := range aSet {
		if !tempMap[v] {
			return false
		}
	}
	return true
}

// []int64 	b切片 是否包含 a切片的所有元素
// a∈b
func AllSubSetInt(aSet, bSet []int64) bool {
	l := new(sync.RWMutex)
	tempMap := make(map[int64]bool, 0)
	for _, v := range bSet {
		l.Lock()
		tempMap[v] = true
		l.Unlock()
	}
	for _, v := range aSet {
		if !tempMap[v] {
			return false
		}
	}
	return true
}

// 相对补集 ComplementRelative
// 若A和B 是集合，则A 在B 中的相对补集是这样一个集合：其元素属于B但不属于A，B - A = { x| x∈B且x∉A}
func SetComplementRelativeInt(aSet, bSet []int) []int {
	l := new(sync.RWMutex)
	tempMap := make(map[int]bool, 0)
	for _, v := range bSet {
		l.Lock()
		tempMap[v] = true
		l.Unlock()
	}
	for _, v := range aSet {
		l.RLock()
		is := tempMap[v]
		l.RUnlock()
		if is {
			l.Lock()
			delete(tempMap, v)
			l.Unlock()
		}
	}
	U := make([]int, 0)
	for k, _ := range tempMap {
		l.RLock()
		U = append(U, k)
		l.RUnlock()
	}
	return U
}
