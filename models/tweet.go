package models

import "time"

type Tweet struct {
	ID        int       `json:"id" db:"id"`
	Message   string    `json:"message" db:"message"`
	Tag       string    `json:"tag" db:"tag"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
