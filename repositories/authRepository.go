package repositories

import (
	"car-management-system/models"
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(conn *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{db: conn}
}

// Create a new user
func (repository *AuthRepository) Create(ctx context.Context, user models.User) (int, error) {
	var id int
	// Added 'role' to insert, assuming default is 'user' or handled by DB default
	err := repository.db.QueryRow(ctx,
		"INSERT INTO users(full_name, email, password_hash) VALUES($1, $2, $3) RETURNING user_id",
		user.Full_name, user.Email, user.PasswordHash).Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, err
}

func (repository *AuthRepository) FindAll(ctx context.Context) ([]models.User, error) {
	sql := "SELECT user_id, full_name, email, password_hash, role FROM users ORDER BY user_id"

	rows, err := repository.db.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]models.User, 0)

	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.User_id,
			&user.Full_name,
			&user.Email,
			&user.PasswordHash,
			&user.Role,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// Update Password (used for Reset Password flow)
func (r *AuthRepository) UpdatePassword(ctx context.Context, userID int, newPasswordHash string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE users SET password_hash = $1 WHERE user_id = $2`,
		newPasswordHash, userID,
	)
	return err
}

// FindByEmail - Returns ONLY ID (Used for Forgot Password check)
func (r *AuthRepository) FindByEmail(ctx context.Context, email string) (int, error) {
	var userID int
	err := r.db.QueryRow(ctx, "SELECT user_id FROM users WHERE email = $1", email).Scan(&userID)
	return userID, err
}

// FindByEmailHash - Returns FULL USER (Used for Login/SignIn)
// FIX: Now returns *models.User instead of int
func (r *AuthRepository) FindByEmailHash(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	// We need ID, Password (to compare), and Role (for JWT)
	query := `SELECT user_id, full_name, email, password_hash, role FROM users WHERE email = $1`

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.User_id,
		&user.Full_name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
	)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetBalance - Added COALESCE to prevent errors if balance is NULL
func (repository *AuthRepository) GetBalance(ctx context.Context, userID int) (float64, error) {
	var balance float64
	// COALESCE(balance, 0) ensures we get 0.0 instead of a Scan error if DB is NULL
	err := repository.db.QueryRow(ctx, "SELECT COALESCE(balance, 0) FROM users WHERE user_id = $1", userID).Scan(&balance)

	if err != nil {
		return 0, err
	}
	return balance, nil
}

// SaveResetToken
func (repository *AuthRepository) SaveResetToken(ctx context.Context, userID int, token string) error {
	expiresAt := time.Now().Add(1 * time.Hour)
	query := `INSERT INTO password_resets (user_id, token, expires_at) VALUES ($1, $2, $3)`
	_, err := repository.db.Exec(ctx, query, userID, token, expiresAt)
	return err
}
