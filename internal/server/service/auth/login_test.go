package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ilya372317/pass-keeper/internal/server/domain"
	"github.com/ilya372317/pass-keeper/internal/server/dto"
	auth_mock "github.com/ilya372317/pass-keeper/internal/server/service/auth/mocks"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
)

func TestService_Login(t *testing.T) {
	type repoSettings struct {
		getUserByIDResult *domain.User
		userEmail         string
		errResult         error
	}
	type tokenGenerateSettings struct {
		result    string
		resultErr error
	}
	type tokenManagerSettings struct {
		generate tokenGenerateSettings
	}
	type want struct {
		err    bool
		result string
	}
	tests := []struct {
		name                 string
		repoSettings         repoSettings
		tokenManagerSettings tokenManagerSettings
		arg                  dto.LoginDTO
		want                 want
	}{
		{
			name: "success simple case",
			repoSettings: repoSettings{
				getUserByIDResult: &domain.User{
					Email:          "email",
					HashedPassword: "pass",
					Salt:           "123",
				},
				userEmail: "email",
				errResult: nil,
			},
			tokenManagerSettings: tokenManagerSettings{
				generate: tokenGenerateSettings{
					result:    "hash123",
					resultErr: nil,
				},
			},
			arg: dto.LoginDTO{
				Email:    "email",
				Password: "pass",
			},
			want: want{
				err:    false,
				result: "hash123",
			},
		},
		{
			name: "failed get user from repo",
			repoSettings: repoSettings{
				getUserByIDResult: &domain.User{},
				userEmail:         "email",
				errResult:         domain.ErrUserNotFound,
			},
			tokenManagerSettings: tokenManagerSettings{},
			arg: dto.LoginDTO{
				Email:    "email",
				Password: "123",
			},
			want: want{
				err:    true,
				result: "",
			},
		},
		{
			name: "given password is incorrect",
			repoSettings: repoSettings{
				getUserByIDResult: &domain.User{
					Email:          "email",
					HashedPassword: "123",
					Salt:           "salt",
				},
				userEmail: "email",
				errResult: nil,
			},
			tokenManagerSettings: tokenManagerSettings{
				generate: tokenGenerateSettings{
					result:    "result",
					resultErr: nil,
				},
			},
			arg: dto.LoginDTO{
				Email:    "email",
				Password: "321",
			},
			want: want{
				err:    true,
				result: "",
			},
		},
		{
			name: "failed generate token",
			repoSettings: repoSettings{
				getUserByIDResult: &domain.User{
					Email:          "email",
					HashedPassword: "123",
					Salt:           "salt",
				},
				userEmail: "email",
				errResult: nil,
			},
			tokenManagerSettings: tokenManagerSettings{
				generate: tokenGenerateSettings{
					result:    "123",
					resultErr: fmt.Errorf("failed generate token"),
				},
			},
			arg: dto.LoginDTO{
				Email:    "email",
				Password: "123",
			},
			want: want{
				err:    true,
				result: "",
			},
		},
	}
	ctrl := gomock.NewController(t)
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := auth_mock.NewMockUserRepository(ctrl)
			tokenManager := auth_mock.NewMockTokenManager(ctrl)
			repoHashedPasswordBytes := sha256.Sum256(
				[]byte(tt.repoSettings.getUserByIDResult.HashedPassword),
			)
			repoHashedPassword := hex.EncodeToString(repoHashedPasswordBytes[:])
			tt.repoSettings.getUserByIDResult.SetHashedPassword(repoHashedPassword)
			hashPasswordBytes := sha256.Sum256([]byte(tt.arg.Password))
			hashPassword := hex.EncodeToString(hashPasswordBytes[:])
			tt.arg.Password = hashPassword

			userRepo.EXPECT().
				GetUserByEmail(ctx, tt.repoSettings.userEmail).
				AnyTimes().
				Return(tt.repoSettings.getUserByIDResult, tt.repoSettings.errResult)
			tokenManager.
				EXPECT().
				Generate(tt.repoSettings.getUserByIDResult).
				AnyTimes().
				Return(tt.tokenManagerSettings.generate.result, tt.tokenManagerSettings.generate.resultErr)

			service := NewAuthService(tokenManager, userRepo)
			got, err := service.Login(ctx, tt.arg)
			if tt.want.err {
				require.Error(t, err)
				return
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.want.result, got)
		})
	}
}
