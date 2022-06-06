package main

import (
	"fmt"
	"go_script/grpc_3/pb"
	"google.golang.org/grpc"
	"net"
	"sync"
	"time"
)

const port = 8082

type server struct{}

// 服务端流模式，拉消息
func (s *server) ServerMod(req *pb.RequestData, res pb.Mod_ServerModServer) error {
	i := 0
	for {
		i++
		//业务代码
		_ = res.Send(&pb.ResponseData{
			Data: fmt.Sprintf("这是发给%s的数据流", req.Data),
		})
		time.Sleep(time.Second * 1)
		if i > 10 {
			break
		}
	}
	return nil
}

// 客户端流模式，推消息
func (s *server) ClientMod(cliStr pb.Mod_ClientModServer) error {
	for {
		//业务代码
		res, err := cliStr.Recv()
		if err != nil {
			fmt.Println("本次客户端流数据发送完了:",err)
			break
		}
		fmt.Println("客户端发来消息：",res.Data)
	}
	return nil
}

// 双向流模式，能推能拉
func (s *server) AllMod(allStr pb.Mod_AllModServer) error {
	wg:=sync.WaitGroup{}
	wg.Add(2)
	//接受客户端消息的协程
	go func() {
		defer wg.Done()
		for  {
			//业务代码
			res, err := allStr.Recv()
			if err != nil {
				fmt.Println("本次客户端流数据发送完了:",err)
				break
			}
			fmt.Println("收到客户端发来消息：",res.Data)
		}
	}()

	//发送消息给客户端的协程
	go func() {
		defer wg.Done()
		i := 0
		for {
			i++
			//业务代码
			_ = allStr.Send(&pb.ResponseData{
				Data: fmt.Sprintf("这是发给客户端的数据流"),
			})
			time.Sleep(time.Second * 1)
			if i > 10 {
				break
			}
		}
	}()
	wg.Wait()
	return nil
}

// 启动
func start() {
	// 1.实例化server
	g := grpc.NewServer()
	// 2.注册逻辑到server中
	pb.RegisterModServer(g, &server{})
	// 3.启动server
	lis, err := net.Listen("tcp", "127.0.0.1:8082")
	if err != nil {
		panic("监听错误:" + err.Error())
	}
	err = g.Serve(lis)
	if err != nil {
		panic("启动错误:" + err.Error())
	}

}

func main() {
	start()
}
