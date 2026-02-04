package repositories

import (
	"car-management-system/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CarRepository struct {
	db *pgxpool.Pool
}

func NewCarRepository(conn *pgxpool.Pool) *CarRepository {
	return &CarRepository{db: conn}
}

func (repostory *CarRepository) Create(ctx context.Context, car models.Car) (int, error) {
	var carId int

	sql := "insert into cars(owner_id, brand, model, year, price, status) values($1, $2, $3, $4, $5, $6) returning car_id"
	err := repostory.db.QueryRow(ctx, sql, car.OwnerId, car.Brand, car.Model, car.Year, car.Price, car.Status).Scan(&carId)
	if err != nil {
		return 0, err
	}
	return carId, err
}
