/*
27. 移除元素
https://leetcode.cn/problems/remove-element/description/
*/

package main

import "log"

func main() {

	nums1 := []int{0, 1, 2, 2, 3, 0, 4, 2}
	val := 2

	removeElement1(nums1, val)

	//rse := removeElement1(nums1, val)
	//log.Println(rse)
	//log.Println(nums1)

	//rse := removeElement2(nums1, val)
	//log.Println(rse)
	//log.Println(nums1)
}

func removeElement1(nums []int, val int) int {
	log.Println(nums)
	i := 0
	for i < len(nums) {
		if nums[i] == val {
			nums = append(nums[:i], nums[i+1:]...)
		} else {
			i++
		}
	}
	log.Println(i)
	return len(nums)
}

func removeElement2(nums []int, val int) int {
	left, right := 0, len(nums)
	for left < right {
		if nums[left] == val {
			nums[left] = nums[right-1]
			right--
		} else {
			left++
		}
	}
	log.Println(nums)
	return left
}
