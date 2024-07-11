package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	metrics "github.com/kiritoxkiriko/go-http-metrics/metrics/prometheus"
	"github.com/kiritoxkiriko/go-http-metrics/middleware"
	"github.com/kiritoxkiriko/go-http-metrics/middleware/std"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	srvAddr     = ":8080"
	metricsAddr = ":8081"
)

func main() {
	// Create our middleware.
	mdlw := middleware.New(middleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{}),
	})

	// Create our router with the metrics middleware.
	r := mux.NewRouter()
	r.Use(std.HandlerProvider("", mdlw))

	// Add paths.
	r.Methods("GET").Path("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	})
	r.Methods("GET").Path("/test1").HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusCreated) })
	r.Methods("GET").Path("/test1/test2").HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusAccepted) })
	r.Methods("GET").Path("/test1/test4").HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNonAuthoritativeInfo) })
	r.Methods("GET").Path("/test2").HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) })
	r.Methods("GET").Path("/test3").HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusResetContent) })

	// Serve our handler.
	go func() {
		log.Printf("server listening at %s", srvAddr)
		if err := http.ListenAndServe(srvAddr, r); err != nil {
			log.Panicf("error while serving: %s", err)
		}
	}()

	// Serve our metrics.
	go func() {
		log.Printf("metrics listening at %s", metricsAddr)
		if err := http.ListenAndServe(metricsAddr, promhttp.Handler()); err != nil {
			log.Panicf("error while serving metrics: %s", err)
		}
	}()

	// Wait until some signal is captured.
	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, syscall.SIGTERM, syscall.SIGINT)
	<-sigC
}
