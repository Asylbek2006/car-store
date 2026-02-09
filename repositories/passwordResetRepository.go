package repositories

import (
	"car-management-system/models"
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PasswordResetRepository struct {
	db *pgxpool.Pool
}

func NewPasswordResetRepository(db *pgxpool.Pool) *PasswordResetRepository {
	return &PasswordResetRepository{db: db}
}

// Create a new token
func (r *PasswordResetRepository) Create(ctx context.Context, userID int, token string) error {
	expiresAt := time.Now().Add(1 * time.Hour) // Token valid for 1 hour
	_, err := r.db.Exec(ctx,
		`INSERT INTO password_resets (user_id, token, expires_at, used) VALUES ($1, $2, $3, false)`,
		userID, token, expiresAt,
	)
	return err
}

// Find token details
func (r *PasswordResetRepository) GetByToken(ctx context.Context, token string) (*models.PasswordReset, error) {
	var pr models.PasswordReset
	err := r.db.QueryRow(ctx,
		`SELECT user_id, token, expires_at, used FROM password_resets WHERE token = $1`,
		token,
	).Scan(&pr.UserID, &pr.Token, &pr.ExpiresAt, &pr.Used)

	if err != nil {
		return nil, err
	}
	return &pr, nil
}

// Mark token as used so it can't be used twice
func (r *PasswordResetRepository) MarkUsed(ctx context.Context, token string) error {
	_, err := r.db.Exec(ctx, `UPDATE password_resets SET used = true WHERE token = $1`, token)
	return err
}
