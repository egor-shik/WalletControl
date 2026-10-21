package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"WalletControl/internal/config"
	deliveryhttp "WalletControl/internal/delivery/http"
	"WalletControl/internal/repository"
	"WalletControl/internal/usecase"
	"WalletControl/internal/worker"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.DatabaseDSN)
	if err != nil {
		log.Fatal(err)
	}
	if err := pool.Ping(ctx); err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	repo := repository.NewBudgetRepository(pool)
	uc := usecase.NewUsecase(repo)

	go worker.StartAnalyticsWorker(ctx, uc, 30*time.Second)

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	handler := deliveryhttp.NewExpenseHandler(uc)
	mux.HandleFunc("POST /expense", handler.AddExpense)
	mux.HandleFunc("GET /analytics", handler.GetAnalytics)

	log.Printf("Server started on port %s", cfg.Port)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
	cancel()
}
