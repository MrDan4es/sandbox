package server

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	pb "github.com/mrdan4es/sandbox/api/fileuploadpb/v1"
)

func FileUploadHandler(c pb.FileUploadServiceClient) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const bufferSize = 3 << 20

		slog.Info("starting to download the file", "buffer", bufferSize)

		start := time.Now()
		r.Body = http.MaxBytesReader(w, r.Body, 10<<30)

		err := r.ParseMultipartForm(32 * 1024)
		if err != nil {
			http.Error(w, "Ошибка при парсинге формы", http.StatusBadRequest)
			slog.Error("parse multipart form", "error", err)
			return
		}

		// Получаем файл из формы
		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Ошибка при получении файла", http.StatusInternalServerError)
			slog.Error("form file", "error", err)
			return
		}
		defer file.Close()

		stream, err := c.UploadUpdateFile(context.Background())
		if err != nil {
			http.Error(w, "Ошибка при открытии стрима", http.StatusInternalServerError)
			slog.Error("stream open file", "error", err)
			return
		}

		if err := stream.Send(&pb.UploadUpdateFileRequest{
			Data: &pb.UploadUpdateFileRequest_FileName{FileName: header.Filename},
		}); err != nil {
			http.Error(w, "Ошибка при отправки названия файла", http.StatusInternalServerError)
			slog.Error("stream send file", "error", err)
			return
		}

		buffer := make([]byte, bufferSize)
		for {
			n, err := file.Read(buffer)
			if err == io.EOF {
				break
			}
			if err != nil {
				http.Error(w, "Ошибка при чтении файла", http.StatusInternalServerError)
				slog.Error("read file", "error", err)
				return
			}

			if sendErr := stream.Send(
				&pb.UploadUpdateFileRequest{
					Data: &pb.UploadUpdateFileRequest_ChunkData{ChunkData: buffer[:n]},
				},
			); sendErr != nil {
				http.Error(w, "Ошибка при записи файла", http.StatusInternalServerError)
				slog.Error("write file", "error", sendErr)
				return
			}
			slog.Info("send chunk", "size", len(buffer[:n]))
		}

		fmt.Fprintf(w, "Файл %s успешно загружен", header.Filename)
		slog.Info("file downloaded successfully", "time", time.Since(start))

		if _, err := stream.CloseAndRecv(); err != nil {
			slog.Error("close stream", "error", err)
		}
	}
}

func FileUploadHandler2(c pb.FileUploadServiceClient) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const bufferSize = 4 << 20

		slog.Info("starting to download the file", "buffer", bufferSize)

		start := time.Now()
		r.Body = http.MaxBytesReader(w, r.Body, 10<<30)

		err := r.ParseMultipartForm(32 * 1024)
		if err != nil {
			http.Error(w, "Ошибка при парсинге формы", http.StatusBadRequest)
			slog.Error("parse multipart form", "error", err)
			return
		}

		// Получаем файл из формы
		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Ошибка при получении файла", http.StatusInternalServerError)
			slog.Error("form file", "error", err)
			return
		}
		defer file.Close()

		stream, err := c.UploadUpdateFile(context.Background())
		if err := stream.Send(&pb.UploadUpdateFileRequest{
			Data: &pb.UploadUpdateFileRequest_FileName{FileName: header.Filename},
		}); err != nil {
			http.Error(w, "Ошибка при открытии стрима", http.StatusInternalServerError)
			slog.Error("stream send file", "error", err)
			return
		}

		if _, err := stream.CloseAndRecv(); err != nil {
			slog.Error("close stream", "error", err)
		}

		// Создаем файл на диске
		dst, err := os.Create(filepath.Join("/home/dmorozov/git/sandbox/test", header.Filename))
		if err != nil {
			http.Error(w, "Ошибка при создании файла", http.StatusInternalServerError)
			slog.Error("create file", "error", err)
			return
		}
		defer dst.Close()

		buffer := make([]byte, bufferSize)
		for {
			n, err := file.Read(buffer)
			if n > 0 {
				if _, writeErr := dst.Write(buffer[:n]); writeErr != nil {
					http.Error(w, "Ошибка при записи файла", http.StatusInternalServerError)
					slog.Error("write file", "error", err)
					return
				}
			}
			if err == io.EOF {
				break
			}
			if err != nil {
				http.Error(w, "Ошибка при чтении файла", http.StatusInternalServerError)
				slog.Error("read file", "error", err)
				return
			}
		}

		fmt.Fprintf(w, "Файл %s успешно загружен", header.Filename)
		slog.Info("file downloaded successfully", "time", time.Since(start))
	}
}
