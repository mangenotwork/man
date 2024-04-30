package main

import (
	"log"
	"time"
)

func main() {

	// case1()

	// case2()

	// case3()

	// case4()

	// case5()

	// case6()

	// case7()

	case8()
}

// case1
// 作用域块变量的作用域
func case1() {
	x := 1
	log.Println(x)
	{
		log.Println(x)
		x := 2 // 隐式声明
		log.Println(x)
	}
	log.Println(x)
}

// case2
// 出现“fallthrough”语句后，它后面只能接下一个子句。
func case2() {
	a := 1
	switch a {
	case 1:
		log.Println("1")
		fallthrough
	case 2:
		log.Println("2")
		fallthrough
	case 3:
		log.Println("3")
	default:
		log.Println("0")
	}
}

// case3
// 数组的声明
func case3() {
	var a1 = [...]int{1, 2, 3, 4, 5}
	log.Println(a1)

	var a2 = [4]int{1, 2, 3, 4}
	log.Println(a2)

	var a3 = [4]int{2: 3, 3: 3} // 指定位置
	log.Println(a3)

	var a4 [20]int
	var a5 = new([20]int)
	log.Println(a4, a5)

}

// case4
// 定义类型别名
type A struct {
	Face int
}
type Aa = A      // 类型别名
func (a A) f()   { log.Println("a", a.Face) }
func (a Aa) f1() { log.Println("a1") } // A, Aa 都实现
func case4() {
	var s A = A{Face: 10}
	s.f()
	s.f1()

	var sa Aa = Aa{Face: 10}
	sa.f()
	sa.f1()
}

// case5
// 函数类型
type funcType func(time.Time) // 定义函数类型
func CTime(start time.Time) {
	log.Println(start)
}
func case5() {
	// 方式1
	var timer funcType = CTime
	timer(time.Now())

	// 方式2
	funcType(CTime)(time.Now())
}

// case6
// 递归
func dg(n int) int {
	if n > 0 {
		log.Println(n)
		return dg(n - 1)
	}
	return n
}
func case6() {
	dg(3)
}

// case7
// 匿名函数
func case7() {
	log.Println(func(x, y int) int { return x + y }(3, 4))

	a := func(x, y int) int { return x + y }
	log.Println(a(6, 6))
	log.Println(a(8, 8))
}

// case8
// 闭包函数
func N(i int) func(a int) int {
	i += 1
	return func(d int) int {
		return i + d
	}
}
func case8() {
	log.Println(N(1)(2))

	f := N(2)
	log.Println(f(3))
}
