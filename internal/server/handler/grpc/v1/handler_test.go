package v1

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/ilya372317/pass-keeper/internal/server/domain"
	"github.com/ilya372317/pass-keeper/internal/server/dto"
	v1_mock "github.com/ilya372317/pass-keeper/internal/server/handler/grpc/v1/mocks"
	pb "github.com/ilya372317/pass-keeper/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

var (
	lis *bufconn.Listener
)

func bufDialer(context.Context, string) (net.Conn, error) {
	conn, err := lis.Dial()
	if err != nil {
		return nil, fmt.Errorf("failed make test connection for grpc: %w", err)
	}

	return conn, nil
}

func setupServer(t *testing.T, server pb.PassServiceServer) *grpc.ClientConn {
	t.Helper()
	ctx := context.Background()
	lis = bufconn.Listen(1024 * 1024)
	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	grpcServer := grpc.NewServer()

	pb.RegisterPassServiceServer(grpcServer, server)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := grpcServer.Serve(lis)
		if err != nil {
			log.Printf("failed to serve gRPC server: %v", err)
		}
	}()

	t.Cleanup(func() {
		grpcServer.GracefulStop()
		err = conn.Close()
		require.NoError(t, err)
		err = lis.Close()
		require.NoError(t, err)
		wg.Wait()
	})

	return conn
}

func TestHandler_Register(t *testing.T) {
	type serviceUser struct {
		createdAT string
		updatedAT string
		email     string
		id        uint
	}
	type serviceConfig struct {
		returnUser serviceUser
		returnErr  error
		argument   dto.RegisterDTO
	}
	type wantUser struct {
		id        int64
		email     string
		createdAT string
		updatedAT string
	}
	type want struct {
		err      bool
		errCode  codes.Code
		wantUser wantUser
	}
	tests := []struct {
		name          string
		argument      *pb.RegisterRequest
		want          want
		serviceConfig serviceConfig
	}{
		{
			name: "success case",
			argument: &pb.RegisterRequest{
				Email:    "1@gmail.com",
				Password: "pass",
			},
			want: want{
				wantUser: wantUser{
					id:        1,
					email:     "1@gmail.com",
					createdAT: "2023-01-01 00:00:00",
					updatedAT: "2023-01-01 00:00:00",
				},
			},
			serviceConfig: serviceConfig{
				returnUser: serviceUser{
					createdAT: "2023-01-01 00:00:00",
					updatedAT: "2023-01-01 00:00:00",
					email:     "1@gmail.com",
					id:        1,
				},
				returnErr: nil,
				argument: dto.RegisterDTO{
					Email:    "1@gmail.com",
					Password: "pass",
				},
			},
		},
		{
			name: "invalid request data case",
			argument: &pb.RegisterRequest{
				Email:    "invalid-email",
				Password: "pass",
			},
			want: want{
				err:      true,
				errCode:  codes.InvalidArgument,
				wantUser: wantUser{},
			},
			serviceConfig: serviceConfig{
				returnUser: serviceUser{
					createdAT: "2023-01-01 00:00:00",
					updatedAT: "2023-01-01 00:00:00",
					email:     "invalid-email",
					id:        1,
				},
				returnErr: nil,
				argument:  dto.RegisterDTO{},
			},
		},
		{
			name: "register failed case",
			argument: &pb.RegisterRequest{
				Email:    "1@gmail.com",
				Password: "pass",
			},
			want: want{
				err:     true,
				errCode: codes.Internal,
			},
			serviceConfig: serviceConfig{
				returnUser: serviceUser{
					createdAT: "2023-01-01 00:00:00",
					updatedAT: "2023-01-01 00:00:00",
					email:     "1@gmail.com",
					id:        1,
				},
				returnErr: fmt.Errorf("failed save user"),
				argument: dto.RegisterDTO{
					Email:    "1@gmail.com",
					Password: "pass",
				},
			},
		},
		{
			name: "user already exists",
			argument: &pb.RegisterRequest{
				Email:    "1@gmail.com",
				Password: "123",
			},
			want: want{
				err:      true,
				errCode:  codes.AlreadyExists,
				wantUser: wantUser{},
			},
			serviceConfig: serviceConfig{
				returnUser: serviceUser{
					createdAT: "2023-01-01 00:00:00",
					updatedAT: "2023-01-01 00:00:00",
					email:     "email",
					id:        1,
				},
				returnErr: domain.ErrUserAlreadyExists,
				argument: dto.RegisterDTO{
					Email:    "1@gmail.com",
					Password: "123",
				},
			},
		},
	}
	ctrl := gomock.NewController(t)
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createdAt, err := time.Parse(time.DateTime, tt.serviceConfig.returnUser.createdAT)
			require.NoError(t, err)
			updatedAt, err := time.Parse(time.DateTime, tt.serviceConfig.returnUser.updatedAT)
			require.NoError(t, err)
			returnServiceUser := &domain.User{
				CreatedAT:      createdAt,
				UpdatedAT:      updatedAt,
				Email:          tt.serviceConfig.returnUser.email,
				HashedPassword: "pass",
				Salt:           "salt",
				ID:             tt.serviceConfig.returnUser.id,
			}
			service := v1_mock.NewMockAuthService(ctrl)
			service.
				EXPECT().
				Register(gomock.Any(), tt.serviceConfig.argument).
				AnyTimes().
				Return(returnServiceUser, tt.serviceConfig.returnErr)
			conn := setupServer(t, New(service))
			client := pb.NewPassServiceClient(conn)

			got, err := client.Register(ctx, tt.argument)
			if tt.want.err {
				require.Error(t, err)
				e, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tt.want.errCode, e.Code())
				return
			} else {
				require.NoError(t, err)
			}

			gotUser := got.User

			assert.Equal(t, tt.want.wantUser.email, gotUser.Email)
			assert.Equal(t, tt.want.wantUser.id, gotUser.Id)
			assert.Equal(t, tt.want.wantUser.createdAT, gotUser.CreatedAt)
			assert.Equal(t, tt.want.wantUser.updatedAT, gotUser.UpdatedAt)
		})
	}
}

func TestHandler_Auth(t *testing.T) {
	type want struct {
		token   string
		err     bool
		errCode codes.Code
	}
	type serviceArgument struct {
		loginDTO dto.LoginDTO
	}
	type serviceResult struct {
		token string
		err   error
	}

	type serviceState struct {
		argument serviceArgument
		result   serviceResult
	}
	tests := []struct {
		name         string
		serviceState serviceState
		argument     *pb.AuthRequest
		want         want
	}{
		{
			name: "success login case",
			serviceState: serviceState{
				argument: serviceArgument{
					loginDTO: dto.LoginDTO{
						Email:    "ilya.otinov@gmail.com",
						Password: "pass",
					},
				},
				result: serviceResult{
					token: "token",
					err:   nil,
				},
			},
			argument: &pb.AuthRequest{
				Email:    "ilya.otinov@gmail.com",
				Password: "pass",
			},
			want: want{
				token: "token",
				err:   false,
			},
		},
		{
			name: "incorrect request data case",
			serviceState: serviceState{
				argument: serviceArgument{
					loginDTO: dto.LoginDTO{
						Email:    "invalid-email",
						Password: "123",
					},
				},
				result: serviceResult{
					token: "",
					err:   nil,
				},
			},
			argument: &pb.AuthRequest{
				Email:    "invalid-email",
				Password: "123",
			},
			want: want{
				err:     true,
				errCode: codes.InvalidArgument,
			},
		},
		{
			name: "user not found case",
			serviceState: serviceState{
				argument: serviceArgument{
					loginDTO: dto.LoginDTO{
						Email:    "1@gmail.com",
						Password: "123",
					},
				},
				result: serviceResult{
					token: "",
					err:   sql.ErrNoRows,
				},
			},
			argument: &pb.AuthRequest{
				Email:    "1@gmail.com",
				Password: "123",
			},
			want: want{
				token:   "",
				err:     true,
				errCode: codes.NotFound,
			},
		},
		{
			name: "internal server error in service",
			serviceState: serviceState{
				argument: serviceArgument{
					loginDTO: dto.LoginDTO{
						Email:    "1@gmail.com",
						Password: "123",
					},
				},
				result: serviceResult{
					token: "",
					err:   fmt.Errorf("internal server error"),
				},
			},
			argument: &pb.AuthRequest{
				Email:    "1@gmail.com",
				Password: "123",
			},
			want: want{
				err:     true,
				errCode: codes.Internal,
			},
		},
		{
			name: "invalid password given",
			serviceState: serviceState{
				argument: serviceArgument{
					loginDTO: dto.LoginDTO{
						Email:    "1@gmail.com",
						Password: "123",
					},
				},
				result: serviceResult{
					token: "",
					err:   domain.ErrInvalidPassword,
				},
			},
			argument: &pb.AuthRequest{
				Email:    "1@gmail.com",
				Password: "123",
			},
			want: want{
				token:   "",
				err:     true,
				errCode: codes.InvalidArgument,
			},
		},
	}
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := v1_mock.NewMockAuthService(ctrl)
			service.
				EXPECT().
				Login(gomock.Any(), tt.serviceState.argument.loginDTO).
				AnyTimes().
				Return(tt.serviceState.result.token, tt.serviceState.result.err)
			server := New(service)

			conn := setupServer(t, server)

			client := pb.NewPassServiceClient(conn)

			got, err := client.Auth(ctx, tt.argument)
			if tt.want.err {
				require.Error(t, err)
				e, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tt.want.errCode, e.Code())
				return
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.want.token, got.AccessToken)
		})
	}
}
