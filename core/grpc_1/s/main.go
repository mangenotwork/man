package main

import (
	"context"
	"fmt"
	"go_script/grpc_1/pb/echo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
	"log"
	"net"
)

func main() {

	lis, err := net.Listen("tcp", ":18881")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(unaryInterceptor), grpc.StreamInterceptor(streamInterceptor))

	// Register EchoServer on the server.
	echo.RegisterEchoServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// unaryInterceptor
// 服务端拦截器
func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	log.Println("[Server] 发起了请求， unaryInterceptor. ")
	log.Println("[Server] ctx = ", ctx)
	log.Println("[Server] req = ", req)
	log.Println("[Server] info = ", info.FullMethod, info.Server)

	// 接收上下文参数
	md, _ := metadata.FromIncomingContext(ctx)
	log.Println(md)

	m, err := handler(ctx, req)
	log.Println("[Server] handler = ", m, err)
	if err != nil {
		log.Println("RPC failed with error %v", err)
	}
	return m, err
}

// wrappedStream wraps around the embedded grpc.ServerStream, and intercepts the RecvMsg and
// SendMsg method call.
type wrappedStream struct {
	grpc.ServerStream
}

func newWrappedStream(s grpc.ServerStream) grpc.ServerStream {
	return &wrappedStream{s}
}

// streamInterceptor
func streamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

	log.Println("[Server] 发起了请求， streamInterceptor. ")
	md, _ := metadata.FromIncomingContext(ss.Context())
	log.Println("[Server] md = ", md)

	err := handler(srv, newWrappedStream(ss))
	if err != nil {
		log.Println("RPC failed with error %v", err)
	}
	return err
}

// grpc 方法服务端实现
type server struct {
	echo.UnimplementedEchoServer
}

func (s *server) UnaryEcho(ctx context.Context, in *echo.EchoRequest) (*echo.EchoResponse, error) {
	fmt.Printf("unary echoing message %q\n", in.Message)
	return &echo.EchoResponse{Message: in.Message}, nil
}

func (s *server) BidirectionalStreamingEcho(stream echo.Echo_BidirectionalStreamingEchoServer) error {
	for {
		in, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			fmt.Printf("server: error receiving from stream: %v\n", err)
			return err
		}
		fmt.Printf("bidi echoing message %q\n", in.Message)
		stream.Send(&echo.EchoResponse{Message: in.Message})
	}
}