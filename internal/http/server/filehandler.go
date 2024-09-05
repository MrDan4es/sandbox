package server

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func FirmwareUploadHandler(w http.ResponseWriter, r *http.Request) {
	if err := uploadHandler(w, r); err != nil {
		slog.Error(err.Error())
	}
	return
	//if err == nil {
	//	return
	//}
	//
	//errStatus := status.Convert(err)
	//errMessage, err := protojson.Marshal(errStatus.Proto())
	//if err != nil {
	//	http.Error(w, "failed to marshal error message to JSON", http.StatusInternalServerError)
	//	return
	//}
	//
	//w.Header().Set("Content-Type", "application/json")
	//http.Error(w, string(errMessage), http.StatusPreconditionFailed)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) error {
	slog.Info("starting uploaded")
	start := time.Now()

	const (
		maxMemory  = 10 << 30
		bufferSize = 1 << 20
	)

	err := r.ParseMultipartForm(maxMemory)

	if err != nil {
		http.Error(w, "failed to parse multipart form", http.StatusInternalServerError)
		return err
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Не удалось получить файл", http.StatusBadRequest)
		return err
	}
	defer file.Close()

	dst, err := os.Create(filepath.Join("/home/mrdan4es/git/bazel/sandbox/test", header.Filename))
	if err != nil {
		http.Error(w, "Не удалось создать файл", http.StatusInternalServerError)
		return err
	}
	defer dst.Close()

	buf := make([]byte, bufferSize)
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			http.Error(w, "Ошибка чтения файла", http.StatusInternalServerError)
			return err
		}
		if n == 0 {
			break
		}

		if _, err := dst.Write(buf[:n]); err != nil {
			http.Error(w, "Ошибка записи файла", http.StatusInternalServerError)
			return err
		}
	}

	fmt.Fprintf(w, "Файл %s успешно загружен!", header.Filename)
	slog.Info("file uploaded", "time", time.Since(start))

	return nil
}

func firmwareUploadHandler(w http.ResponseWriter, r *http.Request) error {
	//if err := s.auth.CheckAccess(r.Context(), "", UploadFirmwareAction); err != nil {
	//	return status.Error(codes.PermissionDenied, "access deny")
	//}

	//currentUserID, err := s.auth.EndUserID(r.Context())
	//if err != nil {
	//	return err
	//}

	reader, err := r.MultipartReader()
	if err != nil {
		return err
	}

	_, err = reader.NextPart()
	if err != nil {
		return err
	}

	//r.Body = http.MaxBytesReader(w, r.Body, MaxFirmwareSize)
	//if err := r.ParseMultipartForm(MaxFirmwareSize); err != nil {
	//	return status.Errorf(codes.InvalidArgument, "firmware is too big: %v", err)
	//}
	//
	//file, _, err := r.FormFile("firmware_bin")
	//if err != nil {
	//	return status.Errorf(codes.InvalidArgument, "failed to extract firmware file from form: %v", err)
	//}
	//
	//firmware, err := s.storeFirmware(r.Context(), currentUserID, file)
	//if err != nil {
	//	return err
	//}
	//
	//firmwareJson, err := protojson.Marshal(firmware)
	//if err != nil {
	//	return status.Errorf(codes.Internal, "failed to marshal firmware metadata: %v", err)
	//}
	//
	//if _, err := w.Write(firmwareJson); err != nil {
	//	return status.Errorf(codes.Internal, "failed to write firmware metadata response: %v", err)
	//}

	return nil
}
