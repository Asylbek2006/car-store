package repositories

import (
	"car-management-system/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(conn *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{db: conn}
}

func (repository *AuthRepository) Create(ctx context.Context, user models.User) (int, error) {
	var id int

	err := repository.db.QueryRow(ctx, "insert into users(full_name, email, password_hash) values($1, $2, $3) returning user_id", user.Full_name, user.Email, user.PasswordHash).Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, err
}

func (repository *AuthRepository) FindAll(ctx context.Context) ([]models.User, error) {
	sql := "select user_id, full_name, email, password_hash from users order by user_id"

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
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (repository *AuthRepository) FindEmail(ctx context.Context, email string) (models.User, error) {
	row := repository.db.QueryRow(ctx, "select user_id, full_name, email, password_hash from users where email = $1", email)

	var user models.User
	err := row.Scan(&user.User_id, &user.Full_name, &user.Email, &user.PasswordHash)
	if err != nil {
		return models.User{}, err
	}
	return user, err
}
