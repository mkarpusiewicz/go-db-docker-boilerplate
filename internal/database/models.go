package database

import (
	"time"
)

type Model struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Entry struct {
	Model
	Name string
}