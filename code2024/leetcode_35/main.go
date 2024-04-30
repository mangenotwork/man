/*
35. 搜索插入位置
https://leetcode.cn/problems/search-insert-position/description/
*/

package main

import "log"

func main() {

	nums := []int{1, 3, 5, 6}
	target := 2

	log.Println(searchInsert1(nums, target))
	log.Println(searchInsert2(nums, target))
	log.Println(searchInsert3(nums, target))
}

func searchInsert1(nums []int, target int) int {
	for i := 0; i < len(nums); i++ {
		if nums[i] >= target {
			return i
		}
	}
	return len(nums)
}

func searchInsert2(nums []int, target int) int {
	if target == 0 {
		return 0
	}
	mn := (len(nums) / 2) - 1
	m := nums[mn]
	log.Println(mn, m)
	if target >= m || mn == 0 {
		for i := mn; i < len(nums); i++ {
			if nums[i] >= target {
				return i
			}
		}
	} else {
		for j := mn; j >= 0; j-- {
			if nums[j] >= target {
				return j
			}
		}
		return 0
	}
	return len(nums)
}

func searchInsert3(nums []int, target int) int {
	n := len(nums)
	left, right := 0, n-1
	ans := n
	for left <= right {
		mid := (right-left)>>1 + left
		if target <= nums[mid] {
			ans = mid
			right = mid - 1
		} else {
			left = mid + 1
		}
	}
	return ans
}
