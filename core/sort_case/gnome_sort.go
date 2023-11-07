package main

/*
侏儒排序（Gnome sort或Stupid sort）是一种排序算法，最初由伊朗计算机工程师Hamid Sarbazi-Azad博士（谢里夫理工大学计算机工程教授）
于2000年提出并被称为“愚蠢排序”（不是 与bogosort混淆），然后由Dick Grune描述并命名为“gnome sort”。 它是一种类似于插入排序的排
序算法，除了将元素移动到适当的位置是通过一系列交换完成的，如冒泡排序。 它在概念上很简单，不需要嵌套循环。
平均或预期的运行时间是O（n2），但如果列表最初几乎排序，则倾向于O（n）。
*/

func GnomeSort(list []int) {
	i := 1
	tmp := 0
	for i < len(list) {
		if list[i] >= list[i-1] {
			i++
		} else {
			tmp = list[i]
			list[i] = list[i-1]
			list[i-1] = tmp
			if i > 1 {
				i--
			}
		}
	}
}
