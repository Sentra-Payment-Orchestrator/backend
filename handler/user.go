package handler

import (
	"context"
	"fmt"

	"github.com/dwikie/sentra-payment-orchestrator/helper"
	"github.com/dwikie/sentra-payment-orchestrator/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserHandler struct {
	pool *pgxpool.Pool
}

type UserHandlerDependencies struct{}

func NewUserHandler(pool *pgxpool.Pool, deps *UserHandlerDependencies) *UserHandler {
	return &UserHandler{pool: pool}
}

func (h *UserHandler) InsertUser(ctx context.Context, payload model.CreateUserpayload) error {
	conn, err := h.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("error while acquiring database connection from pool: %v", err)
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("error while beginning transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	hashedPassword, err := helper.HashPassword(payload.Password)
	if err != nil {
		return fmt.Errorf("error while hashing password: %v", err)
	}

	user := tx.QueryRow(ctx, `
	INSERT INTO users (email, password, status)
	VALUES ($1, $2, $3) RETURNING id
	`, payload.Email, hashedPassword, 0)

	var userID int64
	if err := user.Scan(&userID); err != nil {
		return fmt.Errorf("error while scanning user ID: %v", err)
	}

	_, err = tx.Exec(ctx, `
	INSERT INTO user_profile (user_id, full_name, phone_number)
	VALUES ($1, $2, $3)`,
		userID, payload.FullName, payload.PhoneNumber)
	if err != nil {
		return fmt.Errorf("error while creating user profile: %v", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("error while committing transaction: %v", err)
	}

	return nil
}

func (h *UserHandler) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	if email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}

	user := model.User{}
	conn, err := h.pool.Acquire(ctx)
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

func (h *UserHandler) UpdateLastLogin(ctx context.Context, userID int64) error {
	conn, err := h.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("error while acquiring database connection from pool: %v", err)
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, `
		UPDATE users
		SET last_login_at = NOW(), updated_at = NOW()
		WHERE id = $1
	`, userID)

	if err != nil {
		return fmt.Errorf("error while updating last login: %v", err)
	}

	return nil
}
