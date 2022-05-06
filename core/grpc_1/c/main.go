package main

import (
	"context"
	"fmt"
	"github.com/zituocn/gow/lib/logy"
	"google.golang.org/grpc/metadata"
	"log"
	"time"

	"go_script/grpc_1/pb/echo"
	"google.golang.org/grpc"
)

func main() {

	// Set up a connection to the server.
	conn, err := grpc.Dial("127.0.0.1:18881",
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(unaryInterceptor),
		grpc.WithStreamInterceptor(streamInterceptor))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Make a echo client and send RPCs.
	rgc := echo.NewEchoClient(conn)
	callUnaryEcho(rgc, "hello world")
	//callBidiStreamingEcho(rgc)
}

const fallbackToken = "some-secret-token"

// unaryInterceptor is an example unary interceptor.
// 中间件
func unaryInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

	log.Println("[Clinet] 发起了请求， unaryInterceptor. ")
	log.Println("[Clinet] ctx = ", ctx)
	log.Println("[Clinet] method = ", method)
	log.Println("[Clinet] req = ", req)
	log.Println("[Clinet] reply = ", reply)
	log.Println("[Clinet] cc = ", cc, cc.Target(), cc.GetState()) // 服务地址


	/*
		cc.GetState()  https://github.com/grpc/grpc/blob/master/doc/connectivity-semantics-and-api.md
		From/To	CONNECTING	READY	TRANSIENT_FAILURE	IDLE	SHUTDOWN
		CONNECTING	Incremental progress during connection establishment	All steps needed to establish a connection succeeded	Any failure in any of the steps needed to establish connection	No RPC activity on channel for IDLE_TIMEOUT	Shutdown triggered by application.
		READY		Incremental successful communication on established channel.	Any failure encountered while expecting successful communication on established channel.	No RPC activity on channel for IDLE_TIMEOUT
		OR
		upon receiving a GOAWAY while there are no pending RPCs.	Shutdown triggered by application.
		TRANSIENT_FAILURE	Wait time required to implement (exponential) backoff is over.				Shutdown triggered by application.
		IDLE	Any new RPC activity on the channel				Shutdown triggered by application.
		SHUTDOWN

		为了向 gRPC API（即应用程序代码）的用户隐藏所有这些活动的详细信息，同时公开有关通道状态的有意义的信息，我们使用具有五个状态的状态机，定义如下：

		CONNECTING：通道正在尝试建立连接，并且正在等待名称解析、TCP 连接建立或 TLS 握手中涉及的步骤之一取得进展。这可以用作创建通道时的初始状态。

		READY：通道已通过 TLS 握手（或等效）和协议级（HTTP/2 等）握手成功建立了连接，并且所有后续通信尝试都已成功（或在没有任何已知故障的情况下处于挂起状态）。

		TRANSIENT_FAILURE：出现了一些暂时性故障（例如 TCP 3 次握手超时或套接字错误）。该状态的通道最终会切换到 CONNECTING 状态并尝试再次建立连接。由于重试是通过指数退避完成的，因此无法连接的通道一开始将在此状态下花费很少的时间，但随着尝试反复失败，通道将在此状态下花费越来越多的时间。对于许多非致命故障（例如，由于服务器尚不可用，TCP 连接尝试超时），通道可能会在此状态下花费越来越多的时间。

		IDLE：这是通道甚至没有尝试创建连接的状态，因为缺少新的或挂起的 RPC。在这种状态下可以创建新的 RPC。任何在通道上启动 RPC 的尝试都会将通道推出此状态以进行连接。如果在指定的 IDLE_TIMEOUT 内通道上没有 RPC 活动，即在此期间没有新的或未决（活动）的 RPC，则准备好或正在连接的通道将切换到 IDLE。此外，当没有活动或挂起的 RPC 时接收 GOAWAY 的通道也应该切换到 IDLE 以避免在尝试脱落连接的服务器上的连接过载。我们将使用 300 秒（5 分钟）的默认 IDLE_TIMEOUT。

		SHUTDOWN：此频道已开始关闭。任何新的 RPC 都应该立即失败。挂起的 RPC 可能会继续运行，直到应用程序取消它们。通道可能会进入此状态，因为应用程序明确请求关闭，或者在尝试连接通信期间发生不可恢复的错误。（截至 2015 年 6 月 12 日，不存在归类为不可恢复的已知错误（连接或通信时）。）进入此状态的通道永远不会离开此状态。
	*/

	// 给 ctx 设置参数
	ctx = setCtx(map[string]string{"a":"a", "b":"b", "c":"c"}, cc)

	start := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	end := time.Now()
	// 打印日志
	log.Printf("RPC: %s, start time: %s, end time: %s, err: %v", method, start.Format("Basic"), end.Format(time.RFC3339), err)

	//// 打印日志2
	//clientIp := "" // 获取本机ip
	//if err != nil {
	//	log.Printf("[GRPC] %13v | %v->%v | %s | err = %v",
	//		time.Now().Sub(start),
	//		clientIp,
	//		cc.Target(),
	//		method,
	//		err)
	//} else {
	//	log.Printf("[GRPC] %13v | %v->%v | %s ",
	//		time.Now().Sub(start),
	//		clientIp,
	//		cc.Target(),
	//		method)
	//}

	return err
}


// setCtx set context
func setCtx(kv map[string]string, grpcConn *grpc.ClientConn) context.Context {
	if grpcConn == nil {
		return nil
	}
	value := make([]string,0)
	for k, v := range kv {
		value = append(value, k)
		value = append(value, v)
	}
	return metadata.NewOutgoingContext(context.Background(), metadata.Pairs(value...))
}


// streamInterceptor is an example stream interceptor.
func streamInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {

	log.Println("[Clinet] 发起了请求, streamInterceptor. ")

	s, err := streamer(ctx, desc, cc, method, opts...)
	if err != nil {
		return nil, err
	}
	return newWrappedStream(s), nil
}

func newWrappedStream(s grpc.ClientStream) grpc.ClientStream {
	return &wrappedStream{s}
}

// wrappedStream  wraps around the embedded grpc.ClientStream, and intercepts the RecvMsg and
// SendMsg method call.
type wrappedStream struct {
	grpc.ClientStream
}

func (w *wrappedStream) RecvMsg(m interface{}) error {
	log.Printf("Receive a message (Type: %T) at %v", m, time.Now().Format(time.RFC3339))
	return w.ClientStream.RecvMsg(m)
}

func (w *wrappedStream) SendMsg(m interface{}) error {
	log.Printf("Send a message (Type: %T) at %v", m, time.Now().Format(time.RFC3339))
	return w.ClientStream.SendMsg(m)
}

func callUnaryEcho(client echo.EchoClient, message string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := client.UnaryEcho(ctx, &echo.EchoRequest{Message: message})
	if err != nil {
		log.Fatalf("client.UnaryEcho(_) = _, %v: ", err)
	}
	fmt.Println("UnaryEcho: ", resp.Message)
}