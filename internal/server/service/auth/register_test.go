package auth

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ilya372317/pass-keeper/internal/server/domain"
	"github.com/ilya372317/pass-keeper/internal/server/dto"
	auth_mock "github.com/ilya372317/pass-keeper/internal/server/service/auth/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_Register(t *testing.T) {
	type saveUserConfig struct {
		returnErr error
	}
	type HasUserConfig struct {
		argumentEmail string
		returnResult  bool
		returnErr     error
	}
	type getByEmailConfig struct {
		argument   string
		returnUser *domain.User
		returnErr  error
	}
	type repositoryConfig struct {
		saveUserConfig   saveUserConfig
		getByEmailConfig getByEmailConfig
		hasUserConfig    HasUserConfig
	}
	type argument struct {
		registerDTO dto.RegisterDTO
	}
	type want struct {
		result domain.User
		err    bool
	}
	tests := []struct {
		name             string
		repositoryConfig repositoryConfig
		argument         argument
		want             want
	}{
		{
			name: "success case",
			repositoryConfig: repositoryConfig{
				saveUserConfig: saveUserConfig{returnErr: nil},
				getByEmailConfig: getByEmailConfig{
					argument: "email",
					returnUser: &domain.User{
						Email:          "email",
						HashedPassword: "pass",
						Salt:           "salt",
						ID:             1,
					},
					returnErr: nil,
				},
			},
			argument: argument{
				registerDTO: dto.RegisterDTO{
					Email:    "email",
					Password: "pass",
				},
			},
			want: want{
				result: domain.User{
					Email:          "email",
					HashedPassword: "pass",
					Salt:           "salt",
					ID:             1,
				},
				err: false,
			},
		},
		{
			name: "failed save user case",
			repositoryConfig: repositoryConfig{
				saveUserConfig: saveUserConfig{
					returnErr: fmt.Errorf("failed save user"),
				},
				getByEmailConfig: getByEmailConfig{
					argument:   "email",
					returnUser: &domain.User{},
					returnErr:  nil,
				},
			},
			argument: argument{
				registerDTO: dto.RegisterDTO{
					Email:    "email",
					Password: "123",
				},
			},
			want: want{
				result: domain.User{},
				err:    true,
			},
		},
		{
			name: "failed get user by email case",
			repositoryConfig: repositoryConfig{
				saveUserConfig: saveUserConfig{
					returnErr: nil,
				},
				getByEmailConfig: getByEmailConfig{
					argument:   "email",
					returnUser: nil,
					returnErr:  fmt.Errorf("failed get user by email"),
				},
			},
			argument: argument{
				registerDTO: dto.RegisterDTO{
					Email:    "email",
					Password: "123",
				},
			},
			want: want{
				result: domain.User{},
				err:    true,
			},
		},
		{
			name: "user already exists",
			repositoryConfig: repositoryConfig{
				saveUserConfig: saveUserConfig{
					returnErr: nil,
				},
				getByEmailConfig: getByEmailConfig{
					argument:   "email",
					returnUser: nil,
					returnErr:  nil,
				},
				hasUserConfig: HasUserConfig{
					argumentEmail: "email",
					returnResult:  true,
					returnErr:     nil,
				},
			},
			argument: argument{
				registerDTO: dto.RegisterDTO{
					Email:    "email",
					Password: "pass",
				},
			},
			want: want{
				result: domain.User{},
				err:    true,
			},
		},
		{
			name: "has user return unexpected error",
			repositoryConfig: repositoryConfig{
				saveUserConfig:   saveUserConfig{returnErr: nil},
				getByEmailConfig: getByEmailConfig{argument: "email"},
				hasUserConfig: HasUserConfig{
					argumentEmail: "email",
					returnResult:  false,
					returnErr:     fmt.Errorf("failed check user exists"),
				},
			},
			argument: argument{
				registerDTO: dto.RegisterDTO{
					Email:    "email",
					Password: "pass",
				},
			},
			want: want{
				result: domain.User{},
				err:    true,
			},
		},
	}
	ctrl := gomock.NewController(t)
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := auth_mock.NewMockUserRepository(ctrl)
			repository.
				EXPECT().
				SaveUser(gomock.Any(), gomock.Any()).
				AnyTimes().
				Return(tt.repositoryConfig.saveUserConfig.returnErr)
			repository.
				EXPECT().
				GetUserByEmail(gomock.Any(), tt.repositoryConfig.getByEmailConfig.argument).
				AnyTimes().
				Return(tt.repositoryConfig.getByEmailConfig.returnUser, tt.repositoryConfig.getByEmailConfig.returnErr)
			repository.
				EXPECT().
				HasUser(gomock.Any(), gomock.Any()).
				AnyTimes().
				Return(tt.repositoryConfig.hasUserConfig.returnResult, tt.repositoryConfig.hasUserConfig.returnErr)

			tokenManager := auth_mock.NewMockTokenManager(ctrl)
			service := NewAuthService(tokenManager, repository)

			got, err := service.Register(ctx, tt.argument.registerDTO)
			if tt.want.err {
				require.Error(t, err)
				return
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.want.result.Email, got.Email)
			assert.Equal(t, tt.want.result.HashedPassword, got.HashedPassword)
			assert.Equal(t, tt.want.result.Salt, got.Salt)
			assert.Equal(t, tt.want.result.ID, got.ID)
		})
	}
}
