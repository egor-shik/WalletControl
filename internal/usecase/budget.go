package usecase

import (
    "context"
    "errors"
    "fmt"

    "WalletControl/internal/domain"
)

type TransactionRepository interface {
    AddTransaction(ctx context.Context, amount float64, categoryID int) error
    GetCategoryAnalytics(ctx context.Context) (map[string]float64, error)
}

type CategoryRepository interface {
    CreateCategory(ctx context.Context, name string) (int, error)
    GetCategories(ctx context.Context) ([]domain.Category, error)
}

type BudgetUsecase struct {
    txRepo  TransactionRepository
    catRepo CategoryRepository
}

func NewUsecase(txRepo TransactionRepository, catRepo CategoryRepository) *BudgetUsecase {
    return &BudgetUsecase{txRepo: txRepo, catRepo: catRepo}
}

func (b *BudgetUsecase) AddExpense(ctx context.Context, amount float64, categoryID int) error {
    if amount <= 0 {
        return errors.New("amount must be greater than zero")
    }
    if err := b.txRepo.AddTransaction(ctx, amount, categoryID); err != nil {
        return fmt.Errorf("add transaction: %w", err)
    }
    return nil
}

func (b *BudgetUsecase) GetReport(ctx context.Context) (map[string]float64, error) {
    return b.txRepo.GetCategoryAnalytics(ctx)
}

func (b *BudgetUsecase) CreateCategory(ctx context.Context, name string) (int, error) {
    if name == "" {
        return 0, errors.New("category name is required")
    }
    return b.catRepo.CreateCategory(ctx, name)
}

func (b *BudgetUsecase) ListCategories(ctx context.Context) ([]domain.Category, error) {
    return b.catRepo.GetCategories(ctx)
}