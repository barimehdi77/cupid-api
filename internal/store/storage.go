package store

import "github.com/barimehdi77/cupid-api/internal/database"

type Storage struct {
}

func NewStorage(db *database.DB) *Storage {
	return &Storage{}
}
