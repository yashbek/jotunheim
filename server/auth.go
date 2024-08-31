package server

import (
	"context"


	"google.golang.org/grpc"

)

// const profileHashKey = "baloot"

func Auth() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		// md, ok := metadata.FromIncomingContext(ctx)

		// if !ok {
		// 	return nil, errors.New("bad context")
		// }
		
		// if profileToken := md[profileHashKey][0]; profileToken != "" {
		// 	jwt
		// }

		return handler(ctx, req)
	}
}
