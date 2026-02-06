package repositories

import (
	"car-management-system/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentRepository struct {
	db *pgxpool.Pool
}

func NewPaymentRepository(connection *pgxpool.Pool) *PaymentRepository {
	return &PaymentRepository{db: connection}
}

func (repository *PaymentRepository) Create(ctx context.Context, payment models.Payment) error {
	query := `
		insert into payments (user_id, amount, payment_type)
		values ($1, $2, $3)
	`

	_, err := repository.db.Exec(ctx, query,
		payment.UserID,
		payment.Amount,
		payment.PaymentType,
	)

	return err
}

func (repository *PaymentRepository) GetMyPayments(ctx context.Context, userID int) ([]models.Payment, error) {
	query := `
		select payment_id, user_id, amount, payment_type, created_at
		from payments
		where user_id = $1
		order by created_at desc
	`

	rows, err := repository.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []models.Payment

	for rows.Next() {
		var p models.Payment
		err := rows.Scan(
			&p.PaymentID,
			&p.UserID,
			&p.Amount,
			&p.PaymentType,
			&p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}

	return payments, nil
}

func (r *PaymentRepository) GetAll(ctx context.Context) ([]models.Payment, error) {
	rows, err := r.db.Query(ctx,
		`select payment_id, user_id, amount, payment_type, created_at
		 from payments
		 order by created_at desc`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Payment
	for rows.Next() {
		var p models.Payment
		if err := rows.Scan(
			&p.PaymentID,
			&p.UserID,
			&p.Amount,
			&p.PaymentType,
			&p.CreatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}
