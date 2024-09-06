package main

import (
	"bytes"
	"flag"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
)

func main() {
	fileSize := flag.Int64("size", 1024, "file size in megabytes")
	filename := flag.String("filename", "upload.bin", "file name")
	flag.Parse()

	slog.Info("creating file", "size (MB)", *fileSize)
	data := make([]byte, *fileSize*1000*1000)

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	part, err := writer.CreateFormFile("file", *filename)
	if err != nil {
		slog.Error("create form file", "error", err)
		return
	}

	_, err = io.Copy(part, bytes.NewReader(data))
	if err != nil {
		slog.Error("copy file to multipart form", "error", err)
		return
	}

	err = writer.Close()
	if err != nil {
		slog.Error("close form file", "error", err)
		return
	}

	slog.Info("uploading file")
	resp, err := http.Post("http://localhost:8000/v1/update:upload", writer.FormDataContentType(), &requestBody)
	if err != nil {
		slog.Error("post upload request", "error", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("read response body", "error", err)
		return
	}

	slog.Info("successfully", "body", string(body))
}
