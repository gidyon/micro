package mdutil

import (
	"context"

	"google.golang.org/grpc/metadata"
)

// AddFromCtx extracts metadata in ctx and create an outgoing context with the MD attached
func AddFromCtx(ctx context.Context) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		return metadata.NewOutgoingContext(ctx, md)
	}
	return ctx
}

// AddMD adds metadata to the context
func AddMD(ctx context.Context, md metadata.MD) context.Context {
	return metadata.NewOutgoingContext(ctx, md)
}
