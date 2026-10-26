package repository

import (
    "context"
    "fmt"

    "github.com/jackc/pgx/v5/pgxpool"
)

type BudgetRepository struct {
    pool *pgxpool.Pool
}

func NewBudgetRepository(pool *pgxpool.Pool) *BudgetRepository {
    return &BudgetRepository{pool: pool}
}

func (b *BudgetRepository) AddTransaction(ctx context.Context, amount float64, categoryID int) error {
    _, err := b.pool.Exec(ctx, "INSERT INTO transactions (amount, category_id) VALUES ($1, $2)", amount, categoryID)
    if err != nil {
        return fmt.Errorf("add transaction: %w", err)
    }
    return nil
}

func (b *BudgetRepository) GetCategoryAnalytics(ctx context.Context) (map[string]float64, error) {
    query := "SELECT categories.name, SUM(transactions.amount) AS total FROM transactions JOIN categories ON transactions.category_id = categories.id WHERE transactions.created_at >= NOW() - INTERVAL '1 month' GROUP BY categories.name ORDER BY total DESC;"
    rows, err := b.pool.Query(ctx, query)
    if err != nil {
        return nil, fmt.Errorf("transactions analysis: %w", err)
    }
    defer rows.Close()

    result := make(map[string]float64)
    for rows.Next() {
        var name string
        var total float64
        if err := rows.Scan(&name, &total); err != nil {
            return nil, fmt.Errorf("scan: %w", err)
        }
        result[name] = total
    }
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("rows: %w", err)
    }
    return result, nil
}