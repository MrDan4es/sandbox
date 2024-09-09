package main

import (
	"bytes"
	"crypto/tls"
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

	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	slog.Info("uploading file")

	req, err := http.NewRequest("POST", "http://localhost:8080/v1/update:upload", &requestBody)
	if err != nil {
		slog.Error("create upload request", "error", err)
		return
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("x-auth-info", "")

	resp, err := client.Do(req)
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

	slog.Info(
		"successfully",
		"body", string(body),
		"status_code", resp.StatusCode,
		"status", resp.Status,
	)
}
