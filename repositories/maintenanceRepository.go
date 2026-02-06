package repositories

import (
	"car-management-system/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MaintenanceRepository struct {
	db *pgxpool.Pool
}

func NewMaintenanceRepository(conn *pgxpool.Pool) *MaintenanceRepository {
	return &MaintenanceRepository{db: conn}
}

func (repository *MaintenanceRepository) Create(ctx context.Context, m models.MaintenanceRecord) error {
	query := `
		insert into maintenance_records
		(car_id, service_type, service_date, mileage, cost, notes)
		values ($1,$2,$3,$4,$5,$6)
	`
	_, err := repository.db.Exec(ctx, query,
		m.CarID, m.ServiceType, m.ServiceDate, m.Mileage, m.Cost, m.Notes,
	)
	return err
}

func (repository *MaintenanceRepository) GetByCar(ctx context.Context, carID int) ([]models.MaintenanceRecord, error) {
	query := `
		select maintenance_id, car_id, service_type, service_date, mileage, cost, notes
		from maintenance_records
		where car_id=$1
		order by service_date desc
	`

	rows, err := repository.db.Query(ctx, query, carID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.MaintenanceRecord
	for rows.Next() {
		var maintenance models.MaintenanceRecord
		err := rows.Scan(
			&maintenance.MaintenanceID, &maintenance.CarID, &maintenance.ServiceType,
			&maintenance.ServiceDate, &maintenance.Mileage, &maintenance.Cost, &maintenance.Notes,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, maintenance)
	}
	return list, nil
}

func (repository *MaintenanceRepository) Update(ctx context.Context, m models.MaintenanceRecord) error {
	_, err := repository.db.Exec(ctx, `
		update maintenance_records
		set service_type=$1, service_date=$2, mileage=$3, cost=$4, notes=$5
		where maintenance_id=$6`,
		m.ServiceType,
		m.ServiceDate,
		m.Mileage,
		m.Cost,
		m.Notes,
		m.MaintenanceID,
	)
	return err
}

func (repository *MaintenanceRepository) GetByID(ctx context.Context, maintenanceID int) (models.MaintenanceRecord, error) {
	query := `
		select maintenance_id, car_id, service_type, service_date, mileage, cost, notes
		from maintenance_records
		where maintenance_id = $1
	`

	var m models.MaintenanceRecord
	err := repository.db.QueryRow(ctx, query, maintenanceID).Scan(
		&m.MaintenanceID,
		&m.CarID,
		&m.ServiceType,
		&m.ServiceDate,
		&m.Mileage,
		&m.Cost,
		&m.Notes,
	)

	return m, err
}
