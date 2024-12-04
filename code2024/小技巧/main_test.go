package main

import (
	"testing"
)

func Test_Case1(t *testing.T) {

	a := func() bool {
		t.Log("a")
		return false
	}

	b := func() bool {
		t.Log("b")
		return false
	}

	if a() && b() { // 这里a如果是false那么b就不执行了，那么a起到了检查b的作用
		t.Log("逻辑true")
	}

}

func Test_Case2(t *testing.T) {

}
