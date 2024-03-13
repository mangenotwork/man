/*
	https://leetcode.cn/problems/sort-colors/description/

*/

package main

import "log"

func main() {
	nums := []int{2, 0, 2, 1, 1, 0}

	sortColors3(nums)

	log.Println(nums)
}

func sortColors(nums []int) {
	q1 := make([]int, 0)
	q2 := make([]int, 0)
	q3 := make([]int, 0)
	for _, v := range nums {
		if v == 0 {
			q1 = append(q1, v)
		}
		if v == 1 {
			q2 = append(q2, v)
		}
		if v == 2 {
			q3 = append(q3, v)
		}
	}
	q1 = append(q1, q2...)
	q1 = append(q1, q3...)
	nums = q1
}

func sortColors2(nums []int) {
	zero, two := -1, len(nums)
	for i := 0; i < two; {
		if nums[i] == 1 {
			i++
		} else if nums[i] == 2 {
			two -= 1
			nums[i], nums[two] = nums[two], nums[i]
		} else {
			zero += 1
			nums[i], nums[zero] = nums[zero], nums[i]
			i++
		}
	}
}

func sortColors3(nums []int) {
	p0, p2 := 0, len(nums)-1
	for i := 0; i <= p2; i++ {
		for ; i <= p2 && nums[i] == 2; p2-- {
			nums[i], nums[p2] = nums[p2], nums[i]
		}
		if nums[i] == 0 {
			nums[i], nums[p0] = nums[p0], nums[i]
			p0++
		}
	}
}
