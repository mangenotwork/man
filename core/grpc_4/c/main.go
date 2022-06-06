package main

import (
	"context"
	"fmt"
	"go_script/grpc_3/pb"
	"google.golang.org/grpc"
	"sync"
	"time"
)

var rpc pb.ModClient

// 服务端流模式
func serverStreamDemo()  {
	res,err:=rpc.ServerMod(context.Background(),&pb.RequestData{Data: "服务端流模式"})
	if err != nil {
		panic("rpc请求错误："+err.Error())
	}
	for  {
		data,err:=res.Recv() //
		if err != nil {
			fmt.Println("客户端发送完了:",err)
			return
		}
		fmt.Println("客户端返回数据流值:",data.Data)
	}
}

// 客户端流模式
func clientStreamDemo()  {
	cliStr, err := rpc.ClientMod(context.Background())
	if err != nil {
		panic("rpc请求错误：" + err.Error())
	}
	i := 0
	for {
		i++
		_ = cliStr.Send(&pb.RequestData{
			Data: "客户端流模式",
		})
		time.Sleep(time.Second * 1)
		if i > 10 {
			break
		}
	}
}

// 双向流模式
func clientAndServerStreamDemo()  {
	allStr, _ := rpc.AllMod(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(1)

	//接受服务端消息的协程
	go func() {
		defer wg.Done()
		for {
			//业务代码
			res, err := allStr.Recv()
			if err != nil {
				fmt.Println("本次服务端流数据发送完了:", err)
				break
			}
			fmt.Println("收到服务端发来消息：", res.Data)
		}
	}()

	//发送消息给服务端的协程
	go func() {
		defer wg.Done()
		i := 0
		for {
			i++
			//业务代码
			_ = allStr.Send(&pb.RequestData{
				Data: fmt.Sprintf("这是发给服务端的数据流"),
			})
			time.Sleep(time.Second * 1)
			if i > 10 {
				break
			}
		}
	}()
	wg.Wait()
}

// 启动
func start() {
	conn, err := grpc.Dial("127.0.0.1:8082", grpc.WithInsecure())
	if err != nil {
		panic("rpc连接错误：" + err.Error())
	}
	defer conn.Close()
	rpc = pb.NewModClient(conn) //初始化

	serverStreamDemo() //服务端流模式

	clientStreamDemo()  //客户端流模式

	clientAndServerStreamDemo() // 双向流模式
}

func main() {
	start()
}