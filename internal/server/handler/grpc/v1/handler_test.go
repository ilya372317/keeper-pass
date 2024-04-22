package v1

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"sync"
	"testing"

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
