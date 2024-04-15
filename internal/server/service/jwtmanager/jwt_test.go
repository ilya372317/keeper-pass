package jwtmanager

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ilya372317/pass-keeper/internal/server/domain"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
)

func TestJWTManager_Generate(t *testing.T) {
	type args struct {
		userEmail string
		duration  time.Duration
		secretKey string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "success case",
			args: args{
				userEmail: "ilya.otinov@gmail.com",
				duration:  10 * time.Second,
				secretKey: "secret-key",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := New(tt.args.secretKey, tt.args.duration)
			user := &domain.User{
				Email: tt.args.userEmail,
			}

			got, err := j.Generate(user)
			require.NoError(t, err)

			token, err := jwt.ParseWithClaims(got, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
				_, ok := token.Method.(*jwt.SigningMethodHMAC)
				if !ok {
					return nil, fmt.Errorf("invalid sign method: %w", err)
				}

				return []byte(tt.args.secretKey), nil
			})
			require.NoError(t, err)

			claims, ok := token.Claims.(*UserClaims)
			require.True(t, ok)
			assert.Equal(t, tt.args.userEmail, claims.Email)
		})
	}
}

func TestJWTManager_Verify(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
		want    *UserClaims
	}{
		{
			name:    "success case",
			wantErr: false,
			want: &UserClaims{
				RegisteredClaims: jwt.RegisteredClaims{},
				Email:            "email",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			const key = "secret-key"
			j := New(key, time.Second*10)
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, tt.want)
			tokenStr, err := token.SignedString([]byte(key))
			require.NoError(t, err)

			email, err := j.Verify(tokenStr)
			if tt.wantErr {
				require.Error(t, err)
				return
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.want.Email, email)
		})
	}
}
