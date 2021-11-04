package middleware

import (
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func AddPayloadLogging(
	logger *zap.Logger,
) ([]grpc.UnaryServerInterceptor, []grpc.StreamServerInterceptor) {
	// Shared options for the logger, with a custom gRPC code to log level function.
	o := []grpc_zap.Option{
		grpc_zap.WithLevels(codeToLevel),
	}

	// Make sure that log statements internal to gRPC library are logged using the zapLogger as well.
	// grpc_zap.ReplaceGrpcLoggerV2(logger)

	// Add unary interceptors
	unaryInterceptors := []grpc.UnaryServerInterceptor{
		grpc_ctxtags.UnaryServerInterceptor(
			grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor),
		),
		grpc_zap.UnaryServerInterceptor(logger, o...),
	}

	// Add stream interceptors
	streamInterceptors := []grpc.StreamServerInterceptor{
		grpc_ctxtags.StreamServerInterceptor(
			grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor),
		),
		grpc_zap.StreamServerInterceptor(logger, o...),
	}

	return unaryInterceptors, streamInterceptors
}
