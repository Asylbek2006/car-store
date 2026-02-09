package models

type UpdateDocumentRequest struct {
	DocumentType *string `json:"document_type"`
	ExpiryDate   *string `json:"expiry_date"`
	FileURL      *string `json:"file_url"`
}
