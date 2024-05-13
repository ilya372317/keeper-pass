package interceptor

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/ilya372317/pass-keeper/internal/server/domain"
	"github.com/ilya372317/pass-keeper/internal/server/dto"
	interceptormock "github.com/ilya372317/pass-keeper/internal/server/interceptor/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestAuthInterceptor_authorize(t *testing.T) {
	type want struct {
		err     bool
		errCode codes.Code
	}
	type tokenManagerSettings struct {
		err    error
		claims dto.JWTClaimsDTO
	}
	type userRepositorySettings struct {
		emailArgument string
		user          *domain.User
		err           error
	}
	tests := []struct {
		name                   string
		mdContent              map[string]string
		openMethods            []string
		method                 string
		hasMD                  bool
		tokenManagerSettings   tokenManagerSettings
		userRepositorySettings userRepositorySettings
		want                   want
	}{
		{
			name:      "success case",
			mdContent: map[string]string{"authorization": "123"},
			tokenManagerSettings: tokenManagerSettings{
				err: nil,
				claims: dto.JWTClaimsDTO{
					Email: "email",
				},
			},
			userRepositorySettings: userRepositorySettings{
				emailArgument: "email",
				user: &domain.User{
					CreatedAT:      time.Now(),
					UpdatedAT:      time.Now(),
					Email:          "email",
					HashedPassword: "pass",
					Salt:           "salt",
					ID:             1,
				},
				err: nil,
			},
			want: want{
				err: false,
			},
			hasMD: true,
		},
		{
			name:      "md is not provided",
			mdContent: nil,
			tokenManagerSettings: tokenManagerSettings{
				err: nil,
				claims: dto.JWTClaimsDTO{
					Email: "email",
				},
			},
			userRepositorySettings: userRepositorySettings{
				emailArgument: "email",
				user:          &domain.User{},
				err:           nil,
			},
			want: want{
				err:     true,
				errCode: codes.Unauthenticated,
			},
			hasMD: false,
		},
		{
			name:      "missing authorization MD key",
			mdContent: nil,
			hasMD:     true,
			tokenManagerSettings: tokenManagerSettings{
				err: nil,
				claims: dto.JWTClaimsDTO{
					Email: "email",
				},
			},
			userRepositorySettings: userRepositorySettings{
				emailArgument: "email",
				user:          &domain.User{},
				err:           nil,
			},
			want: want{
				err:     true,
				errCode: codes.Unauthenticated,
			},
		},
		{
			name: "token is invalid case",
			mdContent: map[string]string{
				"authorization": "123",
			},
			hasMD: true,
			tokenManagerSettings: tokenManagerSettings{
				err:    fmt.Errorf("token is invalid"),
				claims: dto.JWTClaimsDTO{},
			},
			userRepositorySettings: userRepositorySettings{
				emailArgument: "email",
				user:          &domain.User{},
				err:           nil,
			},
			want: want{
				err:     true,
				errCode: codes.Unauthenticated,
			},
		},
		{
			name: "failed get user from repository",
			mdContent: map[string]string{
				"authorization": "123",
			},
			hasMD: true,
			tokenManagerSettings: tokenManagerSettings{
				err: nil,
				claims: dto.JWTClaimsDTO{
					Email: "email",
				},
			},
			userRepositorySettings: userRepositorySettings{
				emailArgument: "email",
				user:          nil,
				err:           fmt.Errorf("failed get user from repo"),
			},
			want: want{
				err:     true,
				errCode: codes.Internal,
			},
		},
		{
			name:                   "method is open",
			mdContent:              nil,
			openMethods:            []string{"method", "method1"},
			method:                 "method",
			hasMD:                  false,
			tokenManagerSettings:   tokenManagerSettings{},
			userRepositorySettings: userRepositorySettings{},
			want:                   want{},
		},
	}
	ctrl := gomock.NewController(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md := metadata.New(tt.mdContent)
			ctx := metadata.NewIncomingContext(context.Background(), md)
			if !tt.hasMD {
				ctx = context.Background()
			}
			tokenManager := interceptormock.NewMockTokenManager(ctrl)
			tokenManager.
				EXPECT().
				Verify(gomock.Any()).
				AnyTimes().
				Return(tt.tokenManagerSettings.claims, tt.tokenManagerSettings.err)
			userRepository := interceptormock.NewMockUserRepository(ctrl)
			userRepository.
				EXPECT().
				GetUserByEmail(ctx, tt.userRepositorySettings.emailArgument).
				AnyTimes().
				Return(tt.userRepositorySettings.user, tt.userRepositorySettings.err)

			authInterceptor := NewAuthInterceptor(tokenManager, userRepository)
			ctx, err := authInterceptor.authorize(ctx, tt.method, tt.openMethods)
			if tt.want.err {
				require.Error(t, err)
				e, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tt.want.errCode, e.Code())

				return
			} else {
				require.NoError(t, err)
			}

			if len(tt.method) > 0 {
				return
			}

			got, ok := ctx.Value(domain.CtxUserKey{}).(*domain.User)
			wantUser := tt.userRepositorySettings.user
			require.True(t, ok)
			assert.Equal(t, wantUser.Email, got.Email)
			assert.Equal(t, wantUser.Salt, got.Salt)
			assert.Equal(t, wantUser.HashedPassword, got.HashedPassword)
			assert.Equal(t, wantUser.ID, got.ID)
		})
	}
}
