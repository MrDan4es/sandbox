package server

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func NewFileUploadServer() *http.Server {
	r := mux.NewRouter()
	r.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		})
	})
	r.HandleFunc("/v1/update:upload", FirmwareUploadHandler).Methods("POST")

	return &http.Server{
		Handler:      r,
		Addr:         ":8000",
		WriteTimeout: time.Second * 30,
	}
}
