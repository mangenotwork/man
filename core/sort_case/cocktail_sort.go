package main

/*
鸡尾酒排序思路，先从左边开始进行冒泡排序，第一趟冒泡排序完，最大值在的数组的最右端，然后进行第二趟排序，
第二趟排序从右边开始排序，第二趟结束时，最小值在数组最左端，以此类推，每一趟排序完都能将一个在当前数组
（不包括之前排序得到的最大或者最小的数）中最小或者最大的数放在对应的位置。


*/

func CocktailSort(list []int) {
	tmp := 0
	for i := 0; i < len(list)/2; i++ {
		left := 0
		right := len(list) - 1
		for left <= right {
			if list[left] > list[left+1] {
				tmp = list[left]
				list[left] = list[left+1]
				list[left+1] = tmp
			}
			left++
			if list[right-1] > list[right] {
				tmp = list[right-1]
				list[right-1] = list[right]
				list[right] = tmp
			}
			right--
		}
	}
}
