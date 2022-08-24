package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/apus-run/connect-examples/gin/cmd/server/router"
)

func main() {
	h2s := &http2.Server{}
	srv := &http.Server{
		Addr: ":8800",
		Handler: h2c.NewHandler(
			router.New(),
			h2s,
		),
		ReadHeaderTimeout: time.Second,
		ReadTimeout:       5 * time.Minute,
		WriteTimeout:      5 * time.Minute,
		MaxHeaderBytes:    8 * 1024, // 8KiB
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(
		signals,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGINT,
	)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP listen and serve: %v", err)
		}
	}()

	<-signals
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP shutdown: %v", err)
	}
	signal.Notify(
		signals,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGINT,
	)
}
