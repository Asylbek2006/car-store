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

func NewPasswordResetRepository(connection *pgxpool.Pool) *PasswordResetRepository {
	return &PasswordResetRepository{db: connection}
}

func (repository *PasswordResetRepository) Create(ctx context.Context, userID int, token string, expiresAt time.Time) error {
	_, err := repository.db.Exec(ctx,
		`insert into password_resets (user_id, token, expires_at)
		 values ($1,$2,$3)`,
		userID, token, expiresAt,
	)
	return err
}

func (repository *PasswordResetRepository) GetByToken(ctx context.Context, token string) (models.PasswordReset, error) {
	row := repository.db.QueryRow(ctx,
		`select user_id, token, expires_at, used
		 from password_resets
		 where token = $1`,
		token,
	)

	var pr models.PasswordReset
	err := row.Scan(&pr.UserID, &pr.Token, &pr.ExpiresAt, &pr.Used)
	return pr, err
}

func (repository *PasswordResetRepository) MarkUsed(ctx context.Context, token string) error {
	_, err := repository.db.Exec(ctx,
		`update password_resets set used = true where token = $1`,
		token,
	)
	return err
}
