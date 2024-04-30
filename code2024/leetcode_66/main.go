/*
https://leetcode.cn/problems/plus-one/description/
*/
package main

import "log"

func main() {
	// digits := []int{1, 2, 3}
	digits := []int{9, 9}

	log.Println(plusOne(digits))
}

func plusOne(digits []int) []int {
	for i := len(digits) - 1; i >= 0; i-- {
		if digits[i] == 9 {
			digits[i] = 0
			if i == 0 {
				digits = append([]int{1}, digits...)
				break
			}
			continue
		} else {
			digits[i]++
			break
		}
	}
	return digits
}
