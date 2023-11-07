package main

/*
堆排序（Heapsort）是指利用堆这种数据结构所设计的一种排序算法。大顶堆的根节点是二叉树中最大的节点，每次操作提取此最大节点，完成排序。

大顶堆：每个节点的值都大于或等于其子节点的值，在堆排序算法中用于升序排列；
小顶堆：每个节点的值都小于或等于其子节点的值，在堆排序算法中用于降序排列；

算法步骤
创建一个堆 H[0……n-1]；
把堆首（最大值）和堆尾互换；
把堆的尺寸缩小 1；
重复步骤 2，直到堆的尺寸为 1。

*/

func HeapSort(arr []int) {
	i := 0
	tmp := 0
	for i = len(arr)/2 - 1; i >= 0; i-- {
		arr = sift(arr, i, len(arr))
	}
	for i = len(arr) - 1; i >= 1; i-- {
		tmp = arr[0]
		arr[0] = arr[i]
		arr[i] = tmp
		arr = sift(arr, 0, i)
	}
}

func sift(arr []int, i int, arrLen int) []int {
	done := false
	tmp := 0
	maxChild := 0
	for (i*2+1 < arrLen) && (!done) {
		if i*2+1 == arrLen-1 {
			maxChild = i*2 + 1
		} else if arr[i*2+1] > arr[i*2+2] {
			maxChild = i*2 + 1
		} else {
			maxChild = i*2 + 2
		}
		if arr[i] < arr[maxChild] {
			tmp = arr[i]
			arr[i] = arr[maxChild]
			arr[maxChild] = tmp
			i = maxChild
		} else {
			done = true
		}
	}
	return arr
}
