package main

import (
	"context"
	"fmt"
	"go_script/grpc_2/pb/echo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"log"
	"time"
)

const (
	Scheme      = "scheme"
	ServiceName = "service_name"
)

var addrs = []string{"localhost:50051", "localhost:50052", "localhost:50053", "localhost:50054"}

func callUnaryEcho(c echo.EchoClient, message string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.UnaryEcho(ctx, &echo.EchoRequest{Message: message})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	fmt.Println(r.Message)
}

func makeRPCs(cc *grpc.ClientConn, n int) {
	hwc := echo.NewEchoClient(cc)
	for i := 0; i < n; i++ {
		callUnaryEcho(hwc, "this is examples/load_balancing")
	}
}

func main() {
	//// "pick_first" is the default, so there's no need to set the load balancing policy.
	//pickfirstConn, err := grpc.Dial(
	//	fmt.Sprintf("%s:///%s", Scheme, ServiceName),
	//	grpc.WithTransportCredentials(insecure.NewCredentials()),
	//)
	//if err != nil {
	//	log.Fatalf("did not connect: %v", err)
	//}
	//defer pickfirstConn.Close()
	//
	//fmt.Println("--- calling helloworld.Greeter/SayHello with pick_first ---")
	//makeRPCs(pickfirstConn, 10)
	//
	//fmt.Println()

	// Make another ClientConn with round_robin policy.
	roundrobinConn, err := grpc.Dial(
		fmt.Sprintf("%s:///%s", Scheme, ServiceName),
		// https://github.com/grpc/grpc/blob/master/doc/service_config.md
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`), // This sets the initial balancing policy.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer roundrobinConn.Close()

	fmt.Println("--- calling helloworld.Greeter/SayHello with round_robin ---")
	//makeRPCs(roundrobinConn, 20)
	hwc := echo.NewEchoClient(roundrobinConn)
	for i := 0; i < 20; i++ {
		callUnaryEcho(hwc, "this is load_balancing")
	}
}

// Following is an example name resolver implementation. Read the name
// resolution example to learn more about it.

type ResolverBuilder struct{}

func (*ResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &exampleResolver{
		target: target,
		cc:     cc,
		addrsStore: map[string][]string{
			ServiceName: addrs,
		},
	}
	r.start()
	return r, nil
}
func (*ResolverBuilder) Scheme() string { return Scheme }

type exampleResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	addrsStore map[string][]string
}

func (r *exampleResolver) start() {
	addrStrs := r.addrsStore[r.target.Endpoint]
	addrs := make([]resolver.Address, len(addrStrs))
	for i, s := range addrStrs {
		addrs[i] = resolver.Address{Addr: s}
	}
	r.cc.UpdateState(resolver.State{Addresses: addrs})
}
func (*exampleResolver) ResolveNow(o resolver.ResolveNowOptions) {}
func (*exampleResolver) Close()                                  {}

func init() {
	resolver.Register(&ResolverBuilder{})
}

