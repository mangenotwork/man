package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	//http.HandleFunc("/", hello) // 注册自己业务处理的Hander
	//http.HandleFunc("/echo", echo)
	http.Handle("/", http.FileServer(http.Dir("/home/mange/Documents")))
	server := http.Server{Addr: ":8777"}
	if err := server.ListenAndServe(); err != nil { // 监听处理
		fmt.Println("server start failed")
	}

	//// 通过信号量的方式停止服务，如果有一部分请求进行到一半，处理完成再关闭服务器
	//c := make(chan os.Signal, 1)
	//signal.Notify(c, os.Interrupt)
	//s := <-c
	//fmt.Printf("接收信号：%s\n", s)
	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()
	//if err := server.Shutdown(ctx); err != nil {
	//	fmt.Println("server shutdown failed")
	//}
	//fmt.Println("server exit")
}

var i = 0

func hello(w http.ResponseWriter, r *http.Request) {
	if i%2 == 0 {
		time.Sleep(2 * time.Second)
	}
	i++
	//time.Sleep(2 * time.Second)
	fmt.Fprintln(w, fmt.Sprintf("Hello Go! %d", i))
}

func echo(w http.ResponseWriter, r *http.Request) {
	time.Sleep(2 * time.Second)
	fmt.Fprintln(w, "Hello Go echo!")
}
