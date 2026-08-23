package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/11DingKing/lanzhou-water-governance/internal/config"
	"github.com/11DingKing/lanzhou-water-governance/internal/httpapi"
	"github.com/11DingKing/lanzhou-water-governance/internal/repository"
	"github.com/11DingKing/lanzhou-water-governance/internal/service"
	"github.com/11DingKing/lanzhou-water-governance/internal/storage/sqlite"
	"github.com/11DingKing/lanzhou-water-governance/internal/worker"
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	db, err := sqlite.Open(ctx, cfg.DBPath)
	if err != nil {
		logger.Error("open database", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	users := repository.Users{DB: db.SQL}
	monitoringRepo := repository.Monitoring{DB: db.SQL}
	inspections := repository.Inspections{DB: db.SQL}
	collabRepo := repository.Collaboration{DB: db.SQL}
	wasteRepo := repository.Waste{DB: db.SQL}
	projectsRepo := repository.Projects{DB: db.SQL}
	auditRepo := repository.Audit{DB: db.SQL}
	auth := service.NewAuth(users, cfg.SessionTTL)
	monitoring := service.Monitoring{DB: db.SQL, Repo: monitoringRepo, Inspections: inspections, Audit: auditRepo}
	collaboration := service.Collaboration{DB: db.SQL, Repo: collabRepo, Monitoring: monitoringRepo, Audit: auditRepo}
	waste := service.Waste{Repo: wasteRepo, Audit: auditRepo}
	projects := service.Projects{Repo: projectsRepo, Audit: auditRepo}
	router := httpapi.NewRouter(httpapi.Dependencies{DB: db.SQL, Auth: auth, Monitoring: monitoring, Collaboration: collaboration, Waste: waste, Projects: projects, Logger: logger})
	server := &http.Server{Addr: cfg.Addr, Handler: router, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second, WriteTimeout: 15 * time.Second, IdleTimeout: 60 * time.Second}
	reconciler := worker.Reconciler{DB: db.SQL, Idempotency: repository.Idempotency{DB: db.SQL}, Logger: logger}
	runner := worker.Runner{Interval: cfg.WorkerInterval, MaxAttempts: 3, Logger: logger, Jobs: make(chan worker.Job, 4)}
	go func() {
		if err := runner.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
			logger.Error("worker stopped", "error", err)
		}
	}()
	runner.Jobs <- reconciler.Job()
	go func() {
		<-ctx.Done()
		shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelShutdown()
		_ = server.Shutdown(shutdownCtx)
	}()
	logger.Info("server listening", "addr", cfg.Addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("server failed", "error", err)
		os.Exit(1)
	}
}
