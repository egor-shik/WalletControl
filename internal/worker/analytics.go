package worker

import (
    "context"
    "log"
    "time"
    "WalletControl/internal/usecase"
)

func StartAnalyticsWorker(ctx context.Context, uc *usecase.BudgetUsecase, interval time.Duration) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            log.Println("Analytics worker stopped gracefully")
            return
        case <-ticker.C:
            report, err := uc.GetReport(ctx)
            if err != nil {
                log.Printf("Worker error: %v", err)
                continue
            }
            log.Printf("Background analytics: %v", report)
        }
    }
}