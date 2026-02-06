package repositories

import (
	"car-management-system/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type FuelRepository struct {
	db *pgxpool.Pool
}

func NewFuelRepository(connection *pgxpool.Pool) *FuelRepository {
	return &FuelRepository{db: connection}
}

func (repository *FuelRepository) Create(ctx context.Context, fuel models.FuelLog) error {
	query := `
		insert into fuel_logs (car_id, fill_date, liters, price, mileage)
		values ($1,$2,$3,$4,$5)
	`
	_, err := repository.db.Exec(ctx, query,
		fuel.CarID, fuel.FillDate, fuel.Liters, fuel.Price, fuel.Mileage,
	)
	return err
}

func (repository *FuelRepository) GetByCar(ctx context.Context, carID int) ([]models.FuelLog, error) {
	rows, err := repository.db.Query(ctx,
		`select fuel_id, car_id, fill_date, liters, price, mileage
		 from fuel_logs
		 where car_id=$1
		 order by fill_date desc`, carID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.FuelLog
	for rows.Next() {
		var fuel models.FuelLog
		if err := rows.Scan(
			&fuel.FuelID,
			&fuel.CarID,
			&fuel.FillDate,
			&fuel.Liters,
			&fuel.Price,
			&fuel.Mileage,
		); err != nil {
			return nil, err
		}
		list = append(list, fuel)
	}
	return list, nil
}

func (repository *FuelRepository) Update(ctx context.Context, fuel models.FuelLog) error {
	_, err := repository.db.Exec(ctx,
		`update fuel_logs
		 set fill_date=$1, liters=$2, price=$3, mileage=$4
		 where fuel_id=$5`,
		fuel.FillDate,
		fuel.Liters,
		fuel.Price,
		fuel.Mileage,
		fuel.FuelID,
	)
	return err
}

func (repository *FuelRepository) GetByID(ctx context.Context, fuelID int) (models.FuelLog, error) {
	query := `
		select fuel_id, car_id, fill_date, liters, price, mileage
		from fuel_logs
		where fuel_id = $1
	`

	var fuel models.FuelLog
	err := repository.db.QueryRow(ctx, query, fuelID).Scan(
		&fuel.FuelID,
		&fuel.CarID,
		&fuel.FillDate,
		&fuel.Liters,
		&fuel.Price,
		&fuel.Mileage,
	)

	return fuel, err
}
