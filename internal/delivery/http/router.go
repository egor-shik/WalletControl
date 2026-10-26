package http

import (
    "net/http"

    "WalletControl/internal/delivery/http/handler"
    "WalletControl/internal/delivery/http/middleware"
    "WalletControl/internal/usecase"
)

func NewRouter(uc *usecase.BudgetUsecase) http.Handler {
    mux := http.NewServeMux()

    exp := handler.NewExpenseHandler(uc)
    mux.HandleFunc("POST /expense", exp.AddExpense)
    mux.HandleFunc("GET /analytics", exp.GetAnalytics)

    return middleware.Recovery(middleware.Logger(mux))
}