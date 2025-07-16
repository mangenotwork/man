package main

import (
	"fmt"
	"os"
	"reflect"
	"runtime/pprof"
)

// golang学习 《Head First Go语言程序设计》

func main() {
	//case1()

	//case2()

	//case3()

	//case4()

	case5()
}

// 指针
func case1() {
	var a int = 10
	ap := &a
	println(ap)
	println(*ap)
	*ap = 20
	println(ap)
	println(a)
	a = 30
	println(a)
	pf(&a)     // 值被更新了
	println(a) // 31
	pf2(a)     // 值没有被更新
	println(a) // 31
	// 小结: 想在函数没有返回值没有重新分配内存空间下传入参数使用指针类型
	b := float64(*ap)
	fmt.Printf("%T", b)

}

func pf(a *int) {
	*a += 1
	println(a) // 指针指向传入的内存
}

func pf2(a int) {
	a += 1
	println(&a) // 重新分配了内存
}

func case2() {
	type a struct {
		a1 int
		a2 string
	}

	s := &a{1, "a"}

	s1 := "s1"

	fmt.Printf("%#v\n", s)
	fmt.Printf("%#v\n", s1)
}

// 判断变量是否为结构体类型
func isStruct(v interface{}) bool {
	t := reflect.TypeOf(v)

	// 处理指针情况，获取指针指向的元素类型
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t.Kind() == reflect.Struct
}

type Number int

// 值类型  接收参数 n 是不会被修改的
func (n Number) Tow() {
	n *= 2
}

// 指针类型 接收参数 n指针 会被修改
func (n *Number) TowP() {
	*n *= 2
}

func case3() {
	nt1 := Number(10)
	nt1.Tow()
	println(nt1) // 原本值不会被修改 不会翻倍

	nt2 := Number(10)
	nt2.TowP()
	println(nt2) // 传入参数是指针 会翻倍
}

// 注：  为了一致性，你所有的类型函数接受值类型，或者都接受指针类型，但是你应该避免混用的情况

type data1 struct {
	a int
}

// 接受值类型不会被修改，是拷贝不是原值;
func (d data1) set1(a int) {
	d.a = a
}

// 接受指针类型会修改原值,指向原值并修改
func (d *data1) set2(a int) {
	d.a = a
}

func case4() {
	d1 := data1{1}
	d1.set1(2)
	println(d1.a) // 1

	d1.set2(3)
	println(d1.a) // 3
}

// 将程序中的数据隐藏在一部分代码中而对另一部分不可见的方法称为封装

// case5  手动生成 profile 文件（适合短期程序）
// 对于一些短期运行的程序（如脚本、定时任务），没办法通过 HTTP 接口动态采集数据，这时可以手动生成 profile 文件。
func case5() {

	heavyTask := func() {
		// 模拟CPU密集型任务
		sum := 0
		for i := 0; i < 1e8; i++ {
			sum += i
		}
	}

	// 生成CPU profile文件
	cpuFile, _ := os.Create("cpu.pprof")

	defer cpuFile.Close()

	_ = pprof.StartCPUProfile(cpuFile) // 开始CPU采样

	defer pprof.StopCPUProfile() // 程序结束时停止采样

	// 生成堆内存profile文件

	heapFile, _ := os.Create("heap.pprof")

	defer heapFile.Close()

	defer pprof.WriteHeapProfile(heapFile) // 程序结束前写入堆内存数据

	// 你的业务逻辑，比如一段耗时的计算
	heavyTask()

}

//  go tool pprof 工具使用
//  安装 Graphviz 并配置环境变量
// go tool pprof cpu.pprof
// web
// 火焰图go-torch   需要安装 FlameGraph 还有perl
