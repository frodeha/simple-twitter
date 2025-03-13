package models

import "time"

type Tweet struct {
	ID        int64     `json:"id" db:"id"`
	Message   string    `json:"message" db:"message"`
	Tag       string    `json:"tag" db:"tag"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
