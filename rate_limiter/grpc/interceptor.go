package grpc

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"rate_limiter/limiter"
)

// UnaryRateLimitInterceptor 返回一个 gRPC 一元调用的拦截器，用于限流控制。
// 参数 rl 是一个实现了 Limiter 接口的限流器。
func UnaryRateLimitInterceptor(rl limiter.Limiter) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context, // 请求上下文
		req interface{}, // 请求对象
		info *grpc.UnaryServerInfo, // 包含请求方法的元信息
		handler grpc.UnaryHandler, // 实际处理请求的处理器
	) (interface{}, error) {

		// 调用限流器的 Allow 方法判断是否允许该请求
		if !rl.Allow() {
			// 超出限流阈值，返回 ResourceExhausted 错误码
			return nil, status.Error(codes.ResourceExhausted, "rate limit exceeded")
		}

		// 通过限流器检查，调用真正的处理逻辑
		return handler(ctx, req)
	}
}
