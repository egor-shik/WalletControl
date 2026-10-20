package main

import (
    "context"
    "log"
    "net/http"
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

    ctx := context.Background()

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
    handler := deliveryhttp.NewExpenseHandler(uc)
    mux.HandleFunc("POST /expense", handler.AddExpense)
    mux.HandleFunc("GET /analytics", handler.GetAnalytics)

    log.Printf("Server started on port %s", cfg.Port)
    
    log.Fatal(http.ListenAndServe(":"+cfg.Port, mux))
}