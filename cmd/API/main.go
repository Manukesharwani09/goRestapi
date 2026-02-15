package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Manukesharwani09/goRestapi/internal/config"
	"github.com/Manukesharwani09/goRestapi/internal/http/handlers/student"
	"github.com/Manukesharwani09/goRestapi/internal/storage/sqlite"
)

func main() {

	//losd config

	cfg := config.MustLoad()
	// database steup
	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal("failed to initialize storage", err)
	}
	slog.Info("storage initialized", slog.String("path", cfg.StoragePath), slog.String("env", cfg.Env))
	//ruouter
	router := http.NewServeMux()
	router.HandleFunc("POST /api/students", student.New(storage))
	//server setuo
	server := http.Server{
		Addr:    cfg.HTTPServer.Addr,
		Handler: router,
	}
	slog.Info("server started", slog.String("addr", cfg.HTTPServer.Addr))
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("server failed to start", err)
		}

	}()
	<-done

	slog.Info("shutting serevr")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = server.Shutdown(ctx)
	if err != nil {
		slog.Error("server shutdown error", err.Error())
	}
	slog.Info("server gracefully stopped")
}
