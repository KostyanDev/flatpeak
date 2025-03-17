package app

import (
	"app/internal/client"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"app/internal/config"
	"app/internal/service"
	httpServer "app/internal/transport/http"

	_ "app/docs"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func Run() error {
	cfg, err := config.New[config.Config]()
	if err != nil {
		return err
	}

	log := logrus.New()

	carbonClient := client.NewClient(cfg, log)
	trustService := service.New(context.Background(), log, carbonClient)
	handler := httpServer.New(context.Background(), log, trustService)
	router := mux.NewRouter()
	httpServer.RegisterRoutes(router, handler)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.HTTPServer.Host, cfg.HTTPServer.Port),
		Handler: router,
	}
	stopped := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigint
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("HTTP Server Shutdown Error: %v", err)
		}
		close(stopped)
	}()

	log.Printf("Starting HTTP server on %s", fmt.Sprintf("%s:%d", cfg.HTTPServer.Host, cfg.HTTPServer.Port))

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe Error: %v", err)
	}

	<-stopped

	log.Printf("Have a nice day!")

	return nil
}
