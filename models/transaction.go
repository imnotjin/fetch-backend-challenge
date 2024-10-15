package models

import (
	"time"

	"gorm.io/gorm"
)

// Transaction represents a single transaction in the database.
type Transaction struct {
	gorm.Model
	Payer     string    `json:"payer"`
	Points    int       `json:"points"`
	Timestamp time.Time `json:"timestamp"`
}
