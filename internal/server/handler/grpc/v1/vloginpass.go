package v1

import (
	"context"

	pb "github.com/ilya372317/pass-keeper/proto"
	"google.golang.org/grpc/status"
)

func (h *Handler) ShowLoginPass(ctx context.Context, in *pb.ShowLoginPassRequest) (*pb.ShowLoginPassResponse, error) {
	data, err := h.loginPassService.Show(ctx, int(in.Id))
	if err != nil {
		e := checkErr("loginpass", in.Id, err)
		return nil, status.Errorf(e.Code(), e.String())
	}

	return &pb.ShowLoginPassResponse{LoginPass: &pb.LoginPass{
		Id:       int64(data.ID),
		Login:    data.Login,
		Password: data.Password,
		Metadata: &pb.LoginPassMetadata{Url: data.Metadata.URL},
	}}, nil
}
