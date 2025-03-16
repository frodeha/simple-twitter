package models

import (
	"encoding/json"
	"fmt"
	"time"
)

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

func (e *AggregatedTweets) UnmarshalJSON(data []byte) error {
	temp := struct {
		GroupBy    string          `json:"group_by"`
		Aggregates json.RawMessage `json:"aggregates"`
	}{}

	err := json.Unmarshal(data, &temp)
	if err != nil {
		return err
	}

	switch {
	case temp.GroupBy == "year":
		var aggregates []YearlyAggregate
		err := json.Unmarshal(temp.Aggregates, &aggregates)
		if err != nil {
			return err
		}

		for _, aggregate := range aggregates {
			e.Aggregates = append(e.Aggregates, aggregate)
		}
		e.GroupBy = temp.GroupBy

	case temp.GroupBy == "month":
		var aggregates []MonthlyAggregate
		err := json.Unmarshal(temp.Aggregates, &aggregates)
		if err != nil {
			return err
		}

		for _, aggregate := range aggregates {
			e.Aggregates = append(e.Aggregates, aggregate)
		}
		e.GroupBy = temp.GroupBy

	default:
		return fmt.Errorf("invalid aggregate grouping %s", temp.GroupBy)
	}

	return nil
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
