package main

/*
奇偶排序法的思路是在数组中重复两趟扫描。第一趟扫描选择所有的数据项对，a[j]和a[j+1]，j是奇数(j=1, 3, 5……)。
如果它们的关键字的值次序颠倒，就交换它们。第二趟扫描对所有的偶数数据项进行同样的操作(j=2, 4,6……)。重复进行
这样两趟的排序直到数组全部有序。
*/

func OddEvenSort(arr []int) {
	tmp, isSorted := 0, false
	for isSorted == false {
		isSorted = true
		for i := 1; i < len(arr)-1; i = i + 2 {
			if arr[i] > arr[i+1] {
				tmp = arr[i]
				arr[i] = arr[i+1]
				arr[i+1] = tmp

				isSorted = false
			}
		}
		for i := 0; i < len(arr)-1; i = i + 2 {
			if arr[i] > arr[i+1] {
				tmp = arr[i]
				arr[i] = arr[i+1]
				arr[i+1] = tmp

				isSorted = false
			}
		}
	}
}
