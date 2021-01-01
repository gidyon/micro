package middleware

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
)

func customFunc(ctx context.Context, p interface{}) error {
	return fmt.Errorf("recovering from panic: %v", p)
}

// AddRecovery recovers from gRPC panics from handlers
func AddRecovery() ([]grpc.UnaryServerInterceptor, []grpc.StreamServerInterceptor) {
	// Shared option for the logger, with a custom gRPC code to log level function.
	opt := grpc_recovery.WithRecoveryHandlerContext(customFunc)

	// Recovery handlers should typically be last in the chain so that other middleware
	// (e.g. logging) can operate on the recovered state instead of being directly affected by any panic
	return []grpc.UnaryServerInterceptor{
			grpc_recovery.UnaryServerInterceptor(opt),
		}, []grpc.StreamServerInterceptor{
			grpc_recovery.StreamServerInterceptor(opt),
		}
}
