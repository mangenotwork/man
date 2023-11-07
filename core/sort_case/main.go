package main

import "log"

func main() {
	caseList := []int{5, 8, 1, 3, 9, 41, 2, 6, 778, 45, 1, 3, 96, 4, 54, 12, 87, 95, 3, 248, 87, 12, 3, 9, 45}
	log.Println("排序算法 ---> ")

	list := make([]int, len(caseList))

	copy(list, caseList)
	log.Println("___\nlist: \t", list)
	BubbleSort(list)
	log.Println("冒泡排序：\t", list)

	copy(list, caseList)
	log.Println("___\nlist: \t", list)
	CocktailSort(list)
	log.Println("鸡尾酒排序: \t", list)

	copy(list, caseList)
	log.Println("___\nlist: \t", list)
	CombSort(list)
	log.Println("梳排序: \t", list)

	copy(list, caseList)
	log.Println("___\nlist: \t", list)
	CountingSort(list)
	log.Println("计数排序 : \t", list)

	copy(list, caseList)
	log.Println("___\nlist: \t", list)
	GnomeSort(list)
	log.Println("侏儒排序 : \t", list)

	copy(list, caseList)
	log.Println("___\nlist: \t", list)
	HeapSort(list)
	log.Println("堆排序 : \t", list)

	copy(list, caseList)
	log.Println("___\nlist: \t", list)
	InsertionSort(list)
	log.Println("插入排序 : \t", list)

	copy(list, caseList)
	log.Println("___\nlist: \t", list)
	OddEvenSort(list)
	log.Println("奇偶排序 : \t", list)

	copy(list, caseList)
	log.Println("___\nlist: \t", list)
	SelectionSort(list)
	log.Println("选择排序 : \t", list)

	copy(list, caseList)
	log.Println("___\nlist: \t", list)
	ShellSort(list)
	log.Println("希尔排序 : \t", list)

	tree := &btree{nil}
	copy(list, caseList)
	log.Println("___\nlist: \t", list)
	TreeSort(list, tree)
	log.Println("树排序 : \t", list)
}
