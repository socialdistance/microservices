package models

import (
	"time"

	"github.com/google/uuid"
)

type File struct {
	Uuid       uuid.UUID
	FilePath   string
	Type       string
	Extension  string // "ods, csv"
	Created_at time.Time
}
