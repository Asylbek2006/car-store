package repositories

import (
	"car-management-system/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RentalRepository struct {
	db *pgxpool.Pool
}

func NewRentalRepository(connection *pgxpool.Pool) *RentalRepository {
	return &RentalRepository{db: connection}
}

func (repository *RentalRepository) Create(ctx context.Context, rental models.Rental) error {
	query := `
		insert into rentals
		(car_id, owner_id, renter_id, start_date, end_date, daily_price, total_price, status)
		values ($1,$2,$3,$4,$5,$6,$7,$8)
	`

	_, err := repository.db.Exec(ctx, query,
		rental.CarID,
		rental.OwnerID,
		rental.RenterID,
		rental.StartDate,
		rental.EndDate,
		rental.DailyPrice,
		rental.TotalPrice,
		rental.Status,
	)

	return err
}

func (repository *RentalRepository) MarkCarRented(ctx context.Context, carID int) error {
	query := `update cars set status='rented' where car_id=$1`
	_, err := repository.db.Exec(ctx, query, carID)
	return err
}

func (repository *RentalRepository) UpdateStatus(ctx context.Context, rentalID int, status string) error {
	query := `update rentals set status=$1 where rental_id=$2`
	_, err := repository.db.Exec(ctx, query, status, rentalID)
	return err
}

func (repository *RentalRepository) GetByID(ctx context.Context, rentalID int) (models.Rental, error) {
	query := `
		select rental_id, car_id, owner_id, renter_id, start_date, end_date,
		       daily_price, total_price, status
		from rentals
		where rental_id=$1
	`

	var rental models.Rental
	err := repository.db.QueryRow(ctx, query, rentalID).Scan(
		&rental.RentalID,
		&rental.CarID,
		&rental.OwnerID,
		&rental.RenterID,
		&rental.StartDate,
		&rental.EndDate,
		&rental.DailyPrice,
		&rental.TotalPrice,
		&rental.Status,
	)

	return rental, err
}

func (repository *RentalRepository) GetAll(ctx context.Context) ([]models.Rental, error) {
	rows, err := repository.db.Query(ctx,
		`select rental_id, car_id, owner_id, renter_id,
		        start_date, end_date, daily_price, total_price, status
		 from rentals
		 order by rental_id desc`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Rental
	for rows.Next() {
		var rental models.Rental
		if err := rows.Scan(
			&rental.RentalID,
			&rental.CarID,
			&rental.OwnerID,
			&rental.RenterID,
			&rental.StartDate,
			&rental.EndDate,
			&rental.DailyPrice,
			&rental.TotalPrice,
			&rental.Status,
		); err != nil {
			return nil, err
		}
		list = append(list, rental)
	}
	return list, nil
}
