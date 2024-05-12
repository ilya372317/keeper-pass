package interceptor

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type userCredStorage interface {
	GetAccessToken() (string, error)
}

type AuthInterceptor struct {
	userCredStorage userCredStorage
	openMethods     map[string]any
}

func NewAuthInterceptor(
	userCredStorage userCredStorage,
	openMethods map[string]any,
) *AuthInterceptor {
	return &AuthInterceptor{
		userCredStorage: userCredStorage,
		openMethods:     openMethods,
	}
}

func (i *AuthInterceptor) Unary() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		if _, ok := i.openMethods[strings.ToLower(method)]; !ok {
			authCtx, err := i.attachToken(ctx)
			if err != nil {
				return fmt.Errorf("authentication failed: %w", err)
			}
			return invoker(authCtx, method, req, reply, cc, opts...)
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func (i *AuthInterceptor) attachToken(ctx context.Context) (context.Context, error) {
	token, err := i.userCredStorage.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed get access token from storage: %w", err)
	}
	if len(token) == 0 {
		return nil, fmt.Errorf(
			"before access to this command you need to login. use command login [email] [password]",
		)
	}

	return metadata.AppendToOutgoingContext(ctx, "authorization", token), nil
}
