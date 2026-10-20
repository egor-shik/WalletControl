package main

import (
    "context"
    "log"
    "net/http"
    "time"
    
    "github.com/jackc/pgx/v5/pgxpool"
	
    "WalletControl/internal/worker"
	"WalletControl/internal/repository"
    "WalletControl/internal/usecase"
    deliveryhttp "WalletControl/internal/delivery/http"
)

func main() {
    dsn := 
    ctx := context.Background()

    pool, err := pgxpool.New(ctx, dsn)
    if err != nil { 
        log.Fatal(err)
    }
    defer pool.Close()

    repo := repository.NewBudgetRepository(pool)
    uc := usecase.NewUsecase(repo)

    interval := 30 * time.Second
    go worker.StartAnalyticsWorker(ctx, uc, interval)

    mux := http.NewServeMux()
    handler := deliveryhttp.NewExpenseHandler(uc)
    mux.HandleFunc("POST /expense", handler.AddExpense)
    mux.HandleFunc("GET /analytics", handler.GetAnalytics)

    log.Println("Server started on :8080")

    log.Fatal(http.ListenAndServe(":8080", mux))
}