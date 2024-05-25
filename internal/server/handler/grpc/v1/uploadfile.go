package v1

import (
	"errors"
	"fmt"
	"io"

	pb "github.com/ilya372317/pass-keeper/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type streamReader struct {
	stream pb.PassService_UploadServer
	state  *streamReaderState
}

type streamReaderState struct {
	buff []byte
}

func (sr streamReader) Read(p []byte) (int, error) {
	if len(sr.state.buff) == 0 {
		chunk, err := sr.stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return 0, io.EOF
			}
			return 0, fmt.Errorf("failed read grpc stream: %w", err)
		}
		sr.state.buff = chunk.GetContent()
	}
	n := copy(p, sr.state.buff)
	sr.state.buff = sr.state.buff[n:]

	return n, nil
}

func (h *Handler) Upload(stream pb.PassService_UploadServer) error {
	firstChunk, err := stream.Recv()
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "failed get file name from first chunk: %v", err)
	}
	fileName := firstChunk.GetFilename()
	if len(fileName) == 0 {
		return status.Errorf(codes.InvalidArgument, "file name must be not empty")
	}

	reader := streamReader{
		stream: stream,
		state:  &streamReaderState{},
	}

	if err = h.fileService.Upload(stream.Context(), firstChunk.GetFilename(), reader); err != nil {
		return status.Errorf(codes.Internal, "failed upload file by grpc: %v", err)
	}

	if err = stream.SendAndClose(&pb.UploadStatus{
		Success: true,
		Message: "successfully save file on path " + firstChunk.GetFilename(),
	}); err != nil {
		return status.Errorf(codes.Internal, "failed close upload file stream: %v", err)
	}

	return nil
}
