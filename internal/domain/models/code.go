package models

import (
	"github.com/google/uuid"
	"time"
)

type Code struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	Value        string
	ExpirationAt time.Time
}
