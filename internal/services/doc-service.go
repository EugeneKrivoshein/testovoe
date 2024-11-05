package services

import (
	"testovoe/internal/db"
)

type DocumentService struct {
	dbProvider *db.PostgresProvider
}

func NewDocumentService(provider *db.PostgresProvider) *DocumentService {
	return &DocumentService{dbProvider: provider}
}
