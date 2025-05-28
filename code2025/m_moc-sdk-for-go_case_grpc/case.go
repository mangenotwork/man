package m_moc_sdk_for_go_case_grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/keepalive"
	"sync"
	"time"
)

// 示例代码 该例子介绍了客户端连接的示例，使用了 keepalive 保持了连接的健康状态，防止连接超时， 使用 connectivity 根据连接状态进行重连等操作

// google.golang.org/grpc/keepalive 库介绍
/*
google.golang.org/grpc/keepalive 是 Go 语言中 gRPC 框架的一个核心子包，主要用于管理 gRPC 连接的心跳机制（Keepalive），确保长连接的稳定性和可靠性

场景1  长连接维持（分布式系统核心场景）

// 服务端开启心跳检测（检测空闲连接）
server := grpc.NewServer(
    grpc.KeepaliveParams(keepalive.ServerParameters{
        Time:    10 * time.Second, // 每10秒检查一次空闲连接
        Timeout: 5 * time.Second,
    }),
)

场景2  防止 NAT 超时（网络中间件场景）

// 客户端主动发送心跳（应对NAT超时）
conn, err := grpc.Dial(
    "server-address",
    grpc.WithKeepaliveParams(keepalive.ClientParameters{
        Time:                30 * time.Second, // 每30秒发送一次心跳
        PermitWithoutStream: true,             // 允许无流时发送心跳
    }),
)

场景3  连接生命周期管理（资源释放场景）

// 服务端限制连接最大存活时间为1小时
server := grpc.NewServer(
    grpc.KeepaliveParams(keepalive.ServerParameters{
        MaxConnectionAge:      1 * time.Hour,
        MaxConnectionAgeGrace: 5 * time.Minute, // 宽限期内处理完现有请求
    }),
)


*/

// google.golang.org/grpc/connectivity 库介绍
/*

google.golang.org/grpc/connectivity 是 Go 语言中 gRPC 框架的核心子包之一，主要用于管理 gRPC 连接的生命周期状态，提供了一套标准化的连
接状态监测和管理机制。通过该库，开发者可以实时获取连接的当前状态、监听状态变化事件，并基于状态变化执行相应的业务逻辑（如重试、熔断、负载均衡等）

状态查询
state := conn.GetState()
if state == connectivity.Ready {
    // 执行 RPC 请求
}

状态监听
// 通过 Conn.Watch() 方法注册状态变更监听器，实时响应状态变化事件

ctx, cancel := context.WithCancel(context.Background())
defer cancel()
go func() {
    for s := range conn.Watch(ctx) {
        switch s {
        case connectivity.Ready:
            log.Println("连接已就绪")
        case connectivity.TransientFailure:
            log.Println("连接暂时失败，尝试重试")
        }
    }
}()

*/

var (
	mux             sync.Mutex
	connectionCache map[string]*grpc.ClientConn // 存放连接
)

func init() {
	connectionCache = map[string]*grpc.ClientConn{}
}

// 清除连接
func ClearConnectionCache() {
	connectionCache = map[string]*grpc.ClientConn{}
}

// 获取默认配置，用于连接  authorizer 即为连接的认证参数根据实际情况加入
func getDefaultDialOption(authorizer auth.Authorizer) []grpc.DialOption {
	var opts []grpc.DialOption

	// Debug Mode allows us to talk to wssdagent without a proper handshake
	// This means we can debug and test wssdagent without generating certs
	// and having proper tokens

	// Check if debug mode is on
	if ok := isDebugMode(); ok == nil {
		opts = append(opts, grpc.WithInsecure())
	} else {
		opts = append(opts, grpc.WithTransportCredentials(authorizer.WithTransportAuthorization()))
	}

	// 通过 keepalive 保持了连接的健康状态，防止连接超时
	opts = append(opts, grpc.WithKeepaliveParams(
		keepalive.ClientParameters{
			Time:                1 * time.Minute,
			Timeout:             20 * time.Second,
			PermitWithoutStream: true,
		}))

	opts = append(opts, grpc.WithUnaryInterceptor(intercept.NewErrorParsingInterceptor()))

	return opts
}

// 获取连接状态返回连接是否有效
func isValidConnections(conn *grpc.ClientConn) bool {

	switch conn.GetState() {
	case connectivity.TransientFailure:
		fallthrough
	case connectivity.Shutdown:
		return false
	default:
		return true
	}
}

// 获取连接
func getClientConnection(serverAddress *string, authorizer auth.Authorizer) (*grpc.ClientConn, error) {
	mux.Lock()
	defer mux.Unlock()
	// 获取连接地址
	endpoint := getServerEndpoint(serverAddress)

	conn, ok := connectionCache[endpoint]
	if ok {
		// 连接有效就返回
		if isValidConnections(conn) {
			return conn, nil
		}
		// 连接无效就关闭连接
		conn.Close()
	}

	// 创建连接
	opts := getDefaultDialOption(authorizer)
	conn, err := grpc.Dial(endpoint, opts...)
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}

	// 保存连接
	connectionCache[endpoint] = conn

	return conn, nil
}
