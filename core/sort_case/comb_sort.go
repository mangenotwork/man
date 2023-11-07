package main

/*
梳排序的基本思想：梳排序是在冒泡排序上再做优化的。梳排序是待排序列通过增量分为若干个子序列，然后对子序列进行冒泡排序，
然后一步步减少增量，直到增量减到1为止。梳排序的最后一次排序是冒泡排序。
梳排序的增量是根据递减率减小的，递减率的设定影响着梳排序的效率，最有效的递减率为1.3。

*/

func CombSort(list []int) {
	tmp := 0
	arrLen := len(list)
	gap := arrLen
	for gap > 1 {
		gap = gap * 10 / 13 // 最有效的递减率为1.3, int类型不能直接*1.3 要 *10/13
		for i := 0; i+gap < arrLen; i++ {
			if list[i] > list[i+gap] {
				tmp = list[i]
				list[i] = list[i+gap]
				list[i+gap] = tmp
			}
		}
	}
}
