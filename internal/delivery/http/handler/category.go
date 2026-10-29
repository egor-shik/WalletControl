package handler

import (
    "encoding/json"
    "net/http"

    "WalletControl/internal/usecase"
)

type AddCategoryRequest struct {
    Name    string `json:"name"`
}

type CategoryHandler struct {
    uc *usecase.BudgetUsecase
}

func NewCategoryHandler(uc *usecase.BudgetUsecase) *CategoryHandler {
    return &CategoryHandler{uc: uc}
}

func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
    var req AddCategoryRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "decoding error", http.StatusBadRequest)
        return
    }
    id, err := h.uc.CreateCategory(r.Context(), req.Name)
	if err != nil {
    http.Error(w, "create category error", http.StatusBadRequest)
    return
	}
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func (h *CategoryHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
    categories, err := h.uc.ListCategories(r.Context())
    if err != nil {
        http.Error(w, "get categories error", http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(categories)
}