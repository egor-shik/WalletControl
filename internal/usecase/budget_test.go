package usecase

import (
    "context"
    "errors"
    "testing"

	"WalletControl/internal/domain"
)

type mockTxRepo struct {
	addCalled bool
	addErr error
}

type mockCatRepo struct {
	createCalled bool
	createID int 
	createErr error 
}

func (m *mockTxRepo) AddTransaction(ctx context.Context, amount float64, categoryID int) error {
	m.addCalled = true 
	return m.addErr
}

func (m *mockTxRepo) GetCategoryAnalytics(ctx context.Context) (map[string]float64, error) {
	return nil, nil
}

func (m *mockCatRepo) CreateCategory(ctx context.Context, name string) (int, error) {
	m.createCalled = true
	return m.createID, m.createErr
}

func (m *mockCatRepo) GetCategories(ctx context.Context)  ([]domain.Category, error) {
	return nil, nil
}

func TestAddExpense_Success(t *testing.T) {
	mock := &mockTxRepo{addErr: nil}
	uc := NewUsecase(mock, &mockCatRepo{})

	err := uc.AddExpense(context.Background(), 100, 1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !mock.addCalled {
		t.Error("expected AddTransaction to be called")
	}
}

func TestAddExpense_ZeroAmount(t *testing.T) {
    mock := &mockTxRepo{}
    uc := NewUsecase(mock, &mockCatRepo{})
    
    err := uc.AddExpense(context.Background(), 0, 1)
    
    if err == nil {
        t.Fatal("expected error, got nil")
    }
    if mock.addCalled {
        t.Error("AddTransaction should NOT be called for zero amount")
    }
}

func TestAddExpense_RepoError(t *testing.T) {
    repoErr := errors.New("db is down")
    mock := &mockTxRepo{addErr: repoErr}
    uc := NewUsecase(mock, &mockCatRepo{})
    
    err := uc.AddExpense(context.Background(), 100, 1)
    
    if err == nil {
        t.Fatal("expected error, got nil")
    }
    if !mock.addCalled {
        t.Error("AddTransaction should be called")
    }
}

func TestCreateCategory_Success(t *testing.T) {
    mock := &mockCatRepo{createID: 42, createErr: nil}
    uc := NewUsecase(&mockTxRepo{}, mock)

    id, err := uc.CreateCategory(context.Background(), "Еда")

    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if id != 42 {
        t.Errorf("id = %d, want 42", id)
    }
    if !mock.createCalled {
        t.Error("CreateCategory should be called")
    }
}

func TestCreateCategory_EmptyName(t *testing.T) {
    mock := &mockCatRepo{}
    uc := NewUsecase(&mockTxRepo{}, mock)

    id, err := uc.CreateCategory(context.Background(), "")

    if err == nil {
        t.Fatal("expected error, got nil")
    }
    if id != 0 {
        t.Errorf("id = %d, want 0 on error", id)
    }
    if mock.createCalled {
        t.Error("CreateCategory should NOT be called for empty name")
    }
}