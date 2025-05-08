package grpc

import (
	"context"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"rate_limiter/limiter"  // 引入限流器包
	pb "rate_limiter/proto" // 引入 protobuf 生成的接口
)

// 定义 gRPC 服务器结构体，嵌入生成的未实现服务
type server struct {
	pb.UnimplementedLimiterServer
}

// SayHello 实现 Hello RPC 接口的方法
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	// 返回 Hello 消息
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

// Run 启动 gRPC 服务
func Run() {
	// 监听本地 9090 端口
	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("监听端口失败: %v", err)
	}

	// 初始化固定窗口限流器，设置每秒最多允许 10 个请求
	rl := limiter.NewFixedWindow(10, time.Second)

	// 创建 gRPC 服务实例，并注册限流拦截器
	s := grpc.NewServer(grpc.UnaryInterceptor(UnaryRateLimitInterceptor(rl)))

	// 注册自定义服务
	pb.RegisterLimiterServer(s, &server{})

	log.Println("gRPC 服务已启动，监听端口 :9090")
	// 启动 gRPC 服务并处理请求
	if err := s.Serve(lis); err != nil {
		log.Fatalf("启动 gRPC 服务失败: %v", err)
	}
}
