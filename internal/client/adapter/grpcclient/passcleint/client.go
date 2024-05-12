package passcleint

import (
	pb "github.com/ilya372317/pass-keeper/proto"
)

type PassClient struct {
	c pb.PassServiceClient
}

func New(client pb.PassServiceClient) *PassClient {
	return &PassClient{c: client}
}
