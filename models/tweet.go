package models

import "time"

type Tweet struct {
	ID        int64     `json:"id" db:"id"`
	Message   string    `json:"message" db:"message"`
	Tag       string    `json:"tag" db:"tag"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type AggregatedTweets struct {
	GroupBy    string `json:"group_by"`
	Aggregates []any  `json:"aggregates"`
}

type YearlyAggregate struct {
	Year   int `json:"year" db:"year"`
	Tweets int `json:"tweets" db:"tweets"`
}

type MonthlyAggregate struct {
	Year   int `json:"year" db:"year"`
	Month  int `json:"month" db:"month"`
	Tweets int `json:"tweets" db:"tweets"`
}
