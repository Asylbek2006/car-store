package repositories

import (
	"car-management-system/models"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
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

func (repository *CarRepository) GetByID(ctx context.Context, carID int) (models.Car, error) {
	query := `
		select car_id, owner_id, brand, model, year, price, status, created_at
		from cars
		where car_id = $1
	`
	var car models.Car
	err := repository.db.QueryRow(ctx, query, carID).Scan(
		&car.CarId,
		&car.OwnerId,
		&car.Brand,
		&car.Model,
		&car.Year,
		&car.Price,
		&car.Status,
		&car.CreatedAt,
	)

	return car, err
}

func (repository *CarRepository) Update(ctx context.Context, car models.Car) error {
	query := `
		update cars
		set brand=$1, model=$2, year=$3, price=$4, status=$5
		where car_id=$6
	`

	_, err := repository.db.Exec(
		ctx,
		query,
		car.Brand,
		car.Model,
		car.Year,
		car.Price,
		car.Status,
		car.CarId,
	)

	return err
}

func (repository *CarRepository) Delete(ctx context.Context, carID int) error {
	query := `delete from cars where car_id = $1`
	_, err := repository.db.Exec(ctx, query, carID)
	return err
}

func (repository *CarRepository) MarkAvailable(ctx context.Context, carID int) error {
	query := `
		update cars
		set status = 'for_rent'
		where car_id = $1
	`
	_, err := repository.db.Exec(ctx, query, carID)
	return err
}

func (repository *CarRepository) MarkRented(ctx context.Context, carID int) error {
	query := `update cars set status='rented' where car_id=$1`
	_, err := repository.db.Exec(ctx, query, carID)
	return err
}

func (repository *CarRepository) MarkSold(ctx context.Context, carID int) error {
	query := `update cars set status='sold' where car_id=$1`
	_, err := repository.db.Exec(ctx, query, carID)
	return err
}

func (repository *CarRepository) GetAll(ctx context.Context) ([]models.Car, error) {
	rows, err := repository.db.Query(ctx,
		`select car_id, owner_id, brand, model, year, price, status, created_at
		 from cars
		 order by created_at desc`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Car
	for rows.Next() {
		var car models.Car
		if err := rows.Scan(
			&car.CarId,
			&car.OwnerId,
			&car.Brand,
			&car.Model,
			&car.Year,
			&car.Price,
			&car.Status,
			&car.CreatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, car)
	}
	return list, nil
}

func (repository *CarRepository) UpdateFields(ctx context.Context, car models.Car) error {
	_, err := repository.db.Exec(ctx, `
		update cars
		set brand=$1, model=$2, year=$3, price=$4, status=$5
		where car_id=$6`,
		car.Brand, car.Model, car.Year, car.Price, car.Status, car.CarId,
	)
	return err
}

func (repository *CarRepository) GetRecommendedCars(ctx context.Context, carType string) ([]models.Car, error) {
	query := `
		select car_id, owner_id, brand, model, year, price, status, created_at
		from cars
		where status in ('for_sale','for_rent')
		and (
			($1 = 'SUV' and model ilike '%SUV%') or
			($1 = 'Sedan' and model ilike '%Sedan%') or
			($1 = 'Hatchback' and model ilike '%Hatch%') or
			($1 = 'Electric')
		)
		limit 5
	`

	rows, err := repository.db.Query(ctx, query, carType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cars []models.Car
	for rows.Next() {
		var car models.Car
		rows.Scan(
			&car.CarId,
			&car.OwnerId,
			&car.Brand,
			&car.Model,
			&car.Year,
			&car.Price,
			&car.Status,
			&car.CreatedAt,
		)
		cars = append(cars, car)
	}

	return cars, nil
}

func (r *CarRepository) BuyCarTransaction(ctx context.Context, buyerID int, carID int) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// --- ИСПРАВЛЕНИЕ 1: используем car_id вместо id ---
	var carPrice float64
	var ownerID int
	var status string

	// Было: WHERE id = $1
	// Стало: WHERE car_id = $1
	queryCar := `SELECT price, owner_id, status FROM cars WHERE car_id = $1 FOR UPDATE`

	err = tx.QueryRow(ctx, queryCar, carID).Scan(&carPrice, &ownerID, &status)
	if err == pgx.ErrNoRows {
		return errors.New("машина не найдена")
	} else if err != nil {
		return err
	}

	if status != "for_sale" {
		return errors.New("эта машина не продается")
	}
	if ownerID == buyerID {
		return errors.New("нельзя купить свою машину")
	}

	// --- ИСПРАВЛЕНИЕ 2: используем user_id вместо id ---
	var buyerBalance float64

	// Было: WHERE id = $1
	// Стало: WHERE user_id = $1
	queryBuyer := `SELECT balance FROM users WHERE user_id = $1 FOR UPDATE`

	err = tx.QueryRow(ctx, queryBuyer, buyerID).Scan(&buyerBalance)
	if err != nil {
		return fmt.Errorf("ошибка получения данных покупателя: %v", err)
	}

	if buyerBalance < carPrice {
		return errors.New("недостаточно средств")
	}

	// --- ИСПРАВЛЕНИЕ 3: Обновление балансов (user_id) ---

	// Списание у покупателя
	_, err = tx.Exec(ctx, `UPDATE users SET balance = balance - $1 WHERE user_id = $2`, carPrice, buyerID)
	if err != nil {
		return fmt.Errorf("ошибка списания средств: %v", err)
	}

	// Начисление продавцу
	_, err = tx.Exec(ctx, `UPDATE users SET balance = balance + $1 WHERE user_id = $2`, carPrice, ownerID)
	if err != nil {
		return fmt.Errorf("ошибка зачисления средств: %v", err)
	}

	// --- ИСПРАВЛЕНИЕ 4: Обновление машины (car_id) ---
	_, err = tx.Exec(ctx, `UPDATE cars SET owner_id = $1, status = 'owned' WHERE car_id = $2`, buyerID, carID)
	if err != nil {
		return fmt.Errorf("ошибка обновления владельца: %v", err)
	}

	return tx.Commit(ctx)
}
