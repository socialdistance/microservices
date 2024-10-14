package models

import (
	"time"

	"github.com/google/uuid"
)

// TODO
type StoreEvent struct {
	ID        uuid.UUID
	Title     string
	Started   time.Time
	Ended     time.Time
	ServiceID int
}
