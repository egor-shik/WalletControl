package usecase

import (
    "context"
    "fmt"
	"errors"
)

type BudgetRepository interface {
    AddTransaction(ctx context.Context, amount float64, categoryID int) error
    GetCategoryAnalytics(ctx context.Context) (map[string]float64, error)
}

type BudgetUsecase struct {
    repo BudgetRepository
}

func NewUsecase(repo BudgetRepository) *BudgetUsecase {
	return &BudgetUsecase{repo: repo}
}

func (b *BudgetUsecase) AddExpense(ctx context.Context, amount float64, categoryID int) error {
	if amount <= 0 { return errors.New("amount must be greater than zero") }
	err := b.repo.AddTransaction(ctx, amount, categoryID)
	if err != nil {
		return fmt.Errorf("add transaction: %w", err)
	}
	return nil
}

func (b *BudgetUsecase) GetReport(ctx context.Context) (map[string]float64, error) {
	return b.repo.GetCategoryAnalytics(ctx)
}