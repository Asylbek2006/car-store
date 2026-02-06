package repositories

import (
	"car-management-system/models"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DocumentRepository struct {
	db *pgxpool.Pool
}

func NewDocumentRepository(connection *pgxpool.Pool) *DocumentRepository {
	return &DocumentRepository{db: connection}
}

func (repository *DocumentRepository) Create(ctx context.Context, document models.Document) error {
	_, err := repository.db.Exec(ctx,
		`insert into documents (car_id, document_type, expiry_date, file_url)
		 values ($1,$2,$3,$4)`,
		document.CarID, document.DocumentType, document.ExpiryDate, document.FileURL,
	)
	return err
}

func (repository *DocumentRepository) GetByCar(ctx context.Context, carID int) ([]models.Document, error) {
	rows, err := repository.db.Query(ctx,
		`select document_id, car_id, document_type, expiry_date, file_url
		 from documents where car_id=$1 order by document_id desc`, carID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Document
	for rows.Next() {
		var document models.Document
		if err := rows.Scan(
			&document.DocumentID,
			&document.CarID,
			&document.DocumentType,
			&document.ExpiryDate,
			&document.FileURL,
		); err != nil {
			return nil, err
		}
		list = append(list, document)
	}
	return list, nil
}

func (repository *DocumentRepository) Delete(ctx context.Context, documentID int) error {
	cmd, err := repository.db.Exec(ctx,
		`delete from documents where document_id=$1`, documentID)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (repository *DocumentRepository) GetByID(ctx context.Context, documentID int) (models.Document, error) {
	query := `
		select document_id, car_id, document_type, expiry_date, file_url
		from documents
		where document_id = $1
	`

	var document models.Document
	err := repository.db.QueryRow(ctx, query, documentID).Scan(
		&document.DocumentID,
		&document.CarID,
		&document.DocumentType,
		&document.ExpiryDate,
		&document.FileURL,
	)

	return document, err
}

func (repository *DocumentRepository) Update(ctx context.Context, document models.Document) error {
	_, err := repository.db.Exec(ctx, `
		update documents
		set document_type=$1, expiry_date=$2, file_url=$3
		where document_id=$4`,
		document.DocumentType, document.ExpiryDate, document.FileURL, document.DocumentID,
	)
	return err
}
