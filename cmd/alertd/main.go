package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/example/alert-orchestration-engine/internal/adapter/http"
	"github.com/example/alert-orchestration-engine/internal/application"
	"github.com/example/alert-orchestration-engine/internal/config"
	"github.com/example/alert-orchestration-engine/internal/infrastructure/clock"
	"github.com/example/alert-orchestration-engine/internal/infrastructure/logging"
	"github.com/example/alert-orchestration-engine/internal/infrastructure/memory"
	"github.com/example/alert-orchestration-engine/internal/infrastructure/notifier"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	c, err := config.Load(os.Getenv("ALERT_CONFIG"))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	log := logging.New(c.LogLevel)
	store := memory.NewStore()
	ids := &memory.IDs{}
	sender := notifier.New()
	svc := application.NewService(memory.EventRepo{S: store}, memory.AlertRepo{S: store}, memory.RuleRepo{S: store}, memory.SilenceRepo{S: store}, memory.NotificationRepo{S: store}, sender, clock.System{}, ids)
	metrics := &httpadapter.Metrics{}
	h := httpadapter.NewHandler(svc, metrics, func() bool { return true })
	srv := &http.Server{Addr: fmt.Sprintf("%s:%d", c.Host, c.Port), Handler: httpadapter.Timeout(15*time.Second, httpadapter.Chain(h.Routes(), log)), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second, WriteTimeout: 15 * time.Second, IdleTimeout: 60 * time.Second}
	go func() {
		log.Info("alert orchestration service started", "address", srv.Addr, "environment", c.Environment)
		if e := srv.ListenAndServe(); e != nil && !errors.Is(e, http.ErrServerClosed) {
			log.Error("server stopped", "error", e)
			os.Exit(1)
		}
	}()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
	shutdown, cancel := context.WithTimeout(context.Background(), time.Duration(c.ShutdownSeconds)*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdown); err != nil {
		log.Error("graceful shutdown failed", "error", err)
	}
	log.Info("service stopped")
}
