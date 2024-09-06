package server

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/mrdan4es/sandbox/api/fileuploadpb/v1"
)

type FileUploadServer struct {
	pb.UnimplementedFileUploadServiceServer
}

func New() *FileUploadServer {
	return &FileUploadServer{}
}

func (s *FileUploadServer) UploadUpdateFile(stream pb.FileUploadService_UploadUpdateFileServer) error {
	r, err := stream.Recv()
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	slog.Info("starting download file", "filename", r.GetFileName())

	dst, err := os.Create(filepath.Join("/home/dmorozov/git/sandbox/test", r.GetFileName()))
	if err != nil {
		slog.Error("create file", "error", err)
		return status.Error(codes.Internal, err.Error())
	}
	defer dst.Close()

	readBytes := 0
	for {
		r, err := stream.Recv()
		if err == io.EOF {
			slog.Info("download file finished")
			break
		}
		if err != nil {
			return status.Errorf(codes.Aborted, "uploading unexpectedly aborted: %v", err)
		}

		slog.Info("got chunk", "size", len(r.GetChunkData()))

		readBytes += len(r.GetChunkData())

		if _, err := dst.Write(r.GetChunkData()); err != nil {
			slog.Error("write file", "error", err)
			return status.Error(codes.Internal, err.Error())
		}
	}

	return stream.SendAndClose(&pb.UploadUpdateFileResponse{})
}
