package main

import (
	"log"
	"reflect"
)

func main() {
	a1 := func() { log.Println("test...") }
	Decorator(&a1, a1)
	a1()
}

// Decorator 泛型装饰器
func Decorator(decoPtr, fn interface{}) {
	var decoratedFunc, targetFunc reflect.Value
	decoratedFunc = reflect.ValueOf(decoPtr).Elem()
	targetFunc = reflect.ValueOf(fn)
	v := reflect.MakeFunc(targetFunc.Type(),
		func(in []reflect.Value) (out []reflect.Value) {
			log.Println("before")
			out = targetFunc.Call(in)
			log.Println("after")
			return
		})
	decoratedFunc.Set(v)
	return
}
