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
    cat := handler.NewCategoryHandler(uc)
    
    mux.HandleFunc("POST /expense", exp.AddExpense)
    mux.HandleFunc("GET /analytics", exp.GetAnalytics)
    mux.HandleFunc("POST /categories", cat.CreateCategory)
    mux.HandleFunc("GET /categories", cat.GetCategories)

    return middleware.Recovery(middleware.Logger(mux))
}