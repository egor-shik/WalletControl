package repository

import (
    "context"
    "fmt"

    "WalletControl/internal/domain"
)

func (r *BudgetRepository) CreateCategory(ctx context.Context, name string) (int, error) {
    var id int
    query := `INSERT INTO categories (name) VALUES ($1) RETURNING id`
    if err := r.pool.QueryRow(ctx, query, name).Scan(&id); err != nil {
        return 0, fmt.Errorf("create category: %w", err)
    }
    return id, nil
}

func (r *BudgetRepository) GetCategories(ctx context.Context) ([]domain.Category, error) {
    query := `SELECT id, name FROM categories ORDER BY id`
    rows, err := r.pool.Query(ctx, query)
    if err != nil {
        return nil, fmt.Errorf("get categories: %w", err)
    }
    defer rows.Close()

    var categories []domain.Category
    for rows.Next() {
        var c domain.Category
        if err := rows.Scan(&c.ID, &c.Name); err != nil {
            return nil, fmt.Errorf("scan category: %w", err)
        }
        categories = append(categories, c)
    }
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("rows: %w", err)
    }
    return categories, nil
}