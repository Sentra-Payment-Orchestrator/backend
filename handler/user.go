package handler

import (
	"context"
	"fmt"

	"github.com/dwikie/sentra-payment-orchestrator/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserHandler struct {
	Pool *pgxpool.Pool
}

func NewUserHandler(pool *pgxpool.Pool) *UserHandler {
	return &UserHandler{Pool: pool}
}

func (h *UserHandler) getUserByEmail(email string) (*model.User, error) {
	if email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}

	user := model.User{}
	ctx := context.Background()
	conn, err := h.Pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	query := `SELECT id, email, password, status, last_login_at, created_at, updated_at FROM users WHERE email=$1`
	err = conn.QueryRow(ctx, query, email).Scan(
		&user.Id,
		&user.Email,
		&user.Password,
		&user.Status,
		&user.LastLogin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
