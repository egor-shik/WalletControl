package http

import (
	"WalletControl/internal/usecase"
    "net/http"
	"encoding/json"
)

type AddExpenseRequest struct {
    Amount     float64 `json:"amount"`
    CategoryID int     `json:"category_id"`
	}

type ExpenseHandler struct {
	uc *usecase.BudgetUsecase
	}

func NewExpenseHandler(uc *usecase.BudgetUsecase) *ExpenseHandler {
	return &ExpenseHandler{uc: uc}
	}

func (h *ExpenseHandler) AddExpense(w http.ResponseWriter, r *http.Request) {
	var req AddExpenseRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
    http.Error(w, "decoding error", http.StatusBadRequest)
	return 
	}
	err := h.uc.AddExpense(r.Context(), req.Amount, req.CategoryID)
	if err != nil {
		http.Error(w, "add expense error", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})

	}

	func (h *ExpenseHandler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
		data, err := h.uc.GetReport(r.Context())
		if err != nil {
			http.Error(w, "analytics error", http.StatusInternalServerError)
			return
		} 
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}

	
