// Package models models/base.go
package models

import (
	"time"
)

type Base struct {
	ID        uint      `gorm:"primary"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}
