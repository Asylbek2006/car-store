package repositories

import (
	"car-management-system/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SaleRepository struct {
	db *pgxpool.Pool
}

func NewSaleRepository(connection *pgxpool.Pool) *SaleRepository {
	return &SaleRepository{db: connection}
}

func (repository *SaleRepository) Create(ctx context.Context, sale models.Sale) error {
	query := `
		insert into sales (car_id, seller_id, buyer_id, sale_price)
		values ($1, $2, $3, $4)
	`

	_, err := repository.db.Exec(ctx, query,
		sale.CarID,
		sale.SellerID,
		sale.BuyerID,
		sale.SalePrice,
	)

	return err
}

func (repository *SaleRepository) MarkCarSold(ctx context.Context, carID int) error {
	query := `update cars set status='sold' where car_id=$1`
	_, err := repository.db.Exec(ctx, query, carID)
	return err
}

func (repository *SaleRepository) GetAll(ctx context.Context) ([]models.Sale, error) {
	rows, err := repository.db.Query(ctx,
		`select sale_id, car_id, seller_id, buyer_id, sale_price, sale_date
		 from sales
		 order by sale_date desc`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Sale
	for rows.Next() {
		var sale models.Sale
		if err := rows.Scan(
			&sale.SaleID,
			&sale.CarID,
			&sale.SellerID,
			&sale.BuyerID,
			&sale.SalePrice,
			&sale.SaleDate,
		); err != nil {
			return nil, err
		}
		list = append(list, sale)
	}
	return list, nil
}
