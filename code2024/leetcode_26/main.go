/*
26. 删除有序数组中的重复项
https://leetcode.cn/problems/remove-duplicates-from-sorted-array/description/
*/

package main

import "log"

func main() {

	test1 := []int{1, 1, 2}
	test2 := []int{0, 0, 1, 1, 1, 2, 2, 3, 3, 4}

	log.Println(removeDuplicates(test1))
	log.Println(removeDuplicates(test2))

	log.Println(removeDuplicates2(test1))
	log.Println(removeDuplicates2(test2))

}

func removeDuplicates(nums []int) int {
	rse := make([]int, 0)
	nMap := make(map[int]struct{})
	for _, v := range nums {
		if _, ok := nMap[v]; !ok {
			rse = append(rse, v)
			nMap[v] = struct{}{}
		}
	}
	nums = rse
	return len(rse)
}

func removeDuplicates2(nums []int) int {
	n := len(nums)
	if n == 0 {
		return 0
	}
	slow := 1
	for fast := 1; fast < n; fast++ {
		if nums[fast] != nums[fast-1] {
			nums[slow] = nums[fast]
			slow++
		}
	}
	return slow
}
