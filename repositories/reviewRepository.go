package repositories

import (
	"car-management-system/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ReviewRepository struct {
	db *pgxpool.Pool
}

func NewReviewRepository(connection *pgxpool.Pool) *ReviewRepository {
	return &ReviewRepository{db: connection}
}

func (repository *ReviewRepository) Create(ctx context.Context, review models.Review) error {
	_, err := repository.db.Exec(ctx,
		`insert into reviews (car_id, user_id, rating, comment)
		 values ($1,$2,$3,$4)`,
		review.CarID,
		review.UserID,
		review.Rating,
		review.Comment,
	)
	return err
}

func (repository *ReviewRepository) GetByCar(ctx context.Context, carID int) ([]models.Review, error) {
	rows, err := repository.db.Query(ctx,
		`select review_id, car_id, user_id, rating, comment, created_at
		 from reviews
		 where car_id = $1
		 order by created_at desc`, carID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Review
	for rows.Next() {
		var review models.Review
		if err := rows.Scan(
			&review.ReviewID,
			&review.CarID,
			&review.UserID,
			&review.Rating,
			&review.Comment,
			&review.CreatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, review)
	}
	return list, nil
}

func (repository *ReviewRepository) ExistsByUser(ctx context.Context, carID, userID int) (bool, error) {
	var exists bool
	err := repository.db.QueryRow(ctx,
		`select exists (
			select 1 from reviews where car_id=$1 and user_id=$2
		)`, carID, userID).Scan(&exists)
	return exists, err
}

func (repository *ReviewRepository) Update(ctx context.Context, rv models.Review) error {
	_, err := repository.db.Exec(ctx, `
		update reviews
		set rating=$1, comment=$2
		where review_id=$3`,
		rv.Rating, rv.Comment, rv.ReviewID,
	)
	return err
}

func (repository *ReviewRepository) GetByID(ctx context.Context, reviewID int) (models.Review, error) {
	query := `
		select review_id, car_id, user_id, rating, comment, created_at
		from reviews
		where review_id = $1
	`

	var review models.Review
	err := repository.db.QueryRow(ctx, query, reviewID).Scan(
		&review.ReviewID,
		&review.CarID,
		&review.UserID,
		&review.Rating,
		&review.Comment,
		&review.CreatedAt,
	)

	return review, err
}
