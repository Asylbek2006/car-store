package repositories

import (
	"car-management-system/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ExpenseRepository struct {
	db *pgxpool.Pool
}

func NewExpenseRepository(connection *pgxpool.Pool) *ExpenseRepository {
	return &ExpenseRepository{db: connection}
}

func (repository *ExpenseRepository) Create(ctx context.Context, expense models.Expense) error {
	_, err := repository.db.Exec(ctx,
		`insert into expenses (car_id, expense_type, amount, expense_date, description)
		 values ($1,$2,$3,$4,$5)`,
		expense.CarID, expense.ExpenseType, expense.Amount, expense.ExpenseDate, expense.Description,
	)
	return err
}

func (repository *ExpenseRepository) GetByCar(ctx context.Context, carID int) ([]models.Expense, error) {
	rows, err := repository.db.Query(ctx,
		`select expense_id, car_id, expense_type, amount, expense_date, description
		 from expenses
		 where car_id=$1
		 order by expense_date desc`, carID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Expense
	for rows.Next() {
		var expense models.Expense
		if err := rows.Scan(
			&expense.ExpenseID,
			&expense.CarID,
			&expense.ExpenseType,
			&expense.Amount,
			&expense.ExpenseDate,
			&expense.Description,
		); err != nil {
			return nil, err
		}
		list = append(list, expense)
	}
	return list, nil
}

func (repository *ExpenseRepository) GetByID(ctx context.Context, expenseID int) (models.Expense, error) {
	query := `
		select expense_id, car_id, expense_type, amount, expense_date, description
		from expenses
		where expense_id = $1
	`

	var expense models.Expense
	err := repository.db.QueryRow(ctx, query, expenseID).Scan(
		&expense.ExpenseID,
		&expense.CarID,
		&expense.ExpenseType,
		&expense.Amount,
		&expense.ExpenseDate,
		&expense.Description,
	)

	return expense, err
}

func (repository *ExpenseRepository) Update(ctx context.Context, e models.Expense) error {
	_, err := repository.db.Exec(ctx, `
		update expenses
		set expense_type=$1, amount=$2, expense_date=$3, description=$4
		where expense_id=$5`,
		e.ExpenseType, e.Amount, e.ExpenseDate, e.Description, e.ExpenseID,
	)
	return err
}
