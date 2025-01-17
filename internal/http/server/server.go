package server

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	pb "github.com/mrdan4es/sandbox/api/fileuploadpb/v1"
)

func NewFileUploadServer(c pb.FileUploadServiceClient) *http.Server {
	r := mux.NewRouter()
	reg := prometheus.NewRegistry()
	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	r.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		})
	})
	r.HandleFunc("/v1/update:upload", FileUploadHandler(c)).Methods("POST")
	r.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	return &http.Server{
		Handler:      r,
		Addr:         ":8000",
		WriteTimeout: time.Second * 30,
	}
}
