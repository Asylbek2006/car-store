package repositories

import (
	"context"

	"car-management-system/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminRepository struct {
	db *pgxpool.Pool
}

func NewAdminRepository(connection *pgxpool.Pool) *AdminRepository {
	return &AdminRepository{db: connection}
}

func (repository *AdminRepository) GetStats(ctx context.Context) (models.AdminStatsResponse, error) {
	var stats models.AdminStatsResponse

	repository.db.QueryRow(ctx, `select count(*) from users`).Scan(&stats.Users)
	repository.db.QueryRow(ctx, `select count(*) from cars`).Scan(&stats.Cars)
	repository.db.QueryRow(ctx, `select count(*) from cars where status='for_sale'`).Scan(&stats.CarsForSale)
	repository.db.QueryRow(ctx, `select count(*) from cars where status='for_rent'`).Scan(&stats.CarsForRent)
	repository.db.QueryRow(ctx, `select count(*) from sales`).Scan(&stats.TotalSales)
	repository.db.QueryRow(ctx, `select count(*) from rentals`).Scan(&stats.TotalRentals)

	repository.db.QueryRow(ctx, `
		select coalesce(sum(amount),0)
		from payments
	`).Scan(&stats.TotalRevenue)

	row := repository.db.QueryRow(ctx, `
		select c.car_id, c.owner_id, c.brand, c.model, c.year, c.price, c.status, c.created_at
		from cars c
		join rentals r on r.car_id = c.car_id
		group by c.car_id
		order by count(r.rental_id) desc
		limit 1
	`)
	var car models.Car
	if err := row.Scan(
		&car.CarId,
		&car.OwnerId,
		&car.Brand,
		&car.Model,
		&car.Year,
		&car.Price,
		&car.Status,
		&car.CreatedAt,
	); err == nil {
		stats.TopCar = &car
	}

	return stats, nil
}
