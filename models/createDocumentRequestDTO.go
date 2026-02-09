package models

type CreateDocumentRequest struct {
	DocumentType string `json:"document_type" binding:"required"`
	ExpiryDate   string `json:"expiry_date" binding:"required"`
	FileURL      string `json:"file_url" binding:"required"`
}
