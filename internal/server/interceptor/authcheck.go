package interceptor

import (
	"context"
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/server/domain"
	"github.com/ilya372317/pass-keeper/internal/server/dto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type TokenManager interface {
	Verify(accessToken string) (dto.JWTClaimsDTO, error)
}

type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
}

type AuthInterceptor struct {
	tokenManager   TokenManager
	userRepository UserRepository
}

func NewAuthInterceptor(tokenManager TokenManager, repository UserRepository) *AuthInterceptor {
	return &AuthInterceptor{tokenManager: tokenManager, userRepository: repository}
}

func (a *AuthInterceptor) Unary(openMethods []string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		newCtx, err := a.authorize(ctx, info.FullMethod, openMethods)
		if err != nil {
			return nil, err
		}

		return handler(newCtx, req)
	}
}

func (a *AuthInterceptor) Stream(openMethods []string) grpc.StreamServerInterceptor {
	return func(srv any,
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		_, err := a.authorize(ss.Context(), info.FullMethod, openMethods)
		if err != nil {
			return err
		}

		return handler(srv, ss)
	}
}

func (a *AuthInterceptor) authorize(ctx context.Context, method string, openMethods []string) (context.Context, error) {
	for _, openMethod := range openMethods {
		if openMethod == method {
			return ctx, nil
		}
	}

	fmt.Println(method)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md.Get("authorization")
	if len(values) == 0 {
		return ctx, status.Errorf(codes.Unauthenticated, "authorization token not provided")
	}

	token := values[0]

	claims, err := a.tokenManager.Verify(token)
	if err != nil {
		return ctx, status.Errorf(codes.Unauthenticated, "invalid token given: %v", err)
	}

	user, err := a.userRepository.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		return ctx, status.Errorf(codes.Internal, "failed get authenticated user from storage: %v", err)
	}

	newCtx := context.WithValue(ctx, domain.CtxUserKey{}, user)

	return newCtx, nil
}
