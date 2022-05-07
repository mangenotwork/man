package main

import (
	"context"
	"flag"
	"fmt"
	"go_script/grpc_2/pb/echo"
	"go_script/grpc_2/pb/helloworld"
	"google.golang.org/grpc"
	"log"
	"net"
)

var port = flag.Int("port", 50051, "the port to serve on")

// hwServer is used to implement helloworld.GreeterServer.
type hwServer struct {
	helloworld.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *hwServer) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	log.Println("Hello " + in.Name)
	return &helloworld.HelloReply{Message: "Hello " + in.Name}, nil
}


type ecServer struct {
	echo.UnimplementedEchoServer
}

func (s *ecServer) UnaryEcho(ctx context.Context, req *echo.EchoRequest) (*echo.EchoResponse, error) {
	log.Println("Hello " + req.Message)
	return &echo.EchoResponse{Message: req.Message}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Printf("server listening at %v\n", lis.Addr())

	s := grpc.NewServer()

	// Register Greeter on the server.
	helloworld.RegisterGreeterServer(s, &hwServer{})

	// Register RouteGuide on the same server.
	echo.RegisterEchoServer(s, &ecServer{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}