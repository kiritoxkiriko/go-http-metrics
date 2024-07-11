package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	metrics "github.com/kiritoxkiriko/go-http-metrics/metrics/prometheus"
	"github.com/kiritoxkiriko/go-http-metrics/middleware"
	gojimiddleware "github.com/kiritoxkiriko/go-http-metrics/middleware/goji"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"goji.io"
	"goji.io/pat"
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
	mux := goji.NewMux()
	mux.Use(gojimiddleware.Handler("", mdlw))

	mux.HandleFunc(pat.Get("/"), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Hello world"))
	}))

	// Serve our handler.
	go func() {
		log.Printf("server listening at %s", srvAddr)
		if err := http.ListenAndServe(srvAddr, mux); err != nil {
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
