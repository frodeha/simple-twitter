package database

import (
	"context"
	"database/sql"
	"fmt"
	"simple_twitter/models"
	"time"
)

type DB interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

type TwitterDatabase struct {
	db DB
}

func (t TwitterDatabase) CreateTweet(ctx context.Context, message string, tag string) (int64, error) {
	result, err := t.db.ExecContext(
		ctx,
		`
			INSERT INTO Tweets (message, tag)
			VALUES (?, ?)
		`,
		message, tag,
	)

	if err != nil {
		return 0, fmt.Errorf("failed to insert tweet: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get id of created tweet: %w", err)
	}

	return id, nil
}

func (t TwitterDatabase) GetTweet(ctx context.Context, id int64) (models.Tweet, error) {
	var tweet models.Tweet
	err := t.db.GetContext(
		ctx,
		&tweet,
		`
			SELECT id, message, tag, created_at
			FROM Tweets
			WHERE id = ?
		`,
		id,
	)

	if err == sql.ErrNoRows {
		return models.Tweet{}, models.ErrMissingf("found no tweet with id %d", id)
	}

	if err != nil {
		return models.Tweet{}, fmt.Errorf("failed to get tweet: %w", err)
	}

	return tweet, nil
}

func (t TwitterDatabase) ListTweets(ctx context.Context, tag string, offset int, limit int) ([]models.Tweet, error) {
	tweets := []models.Tweet{}
	err := t.db.SelectContext(
		ctx,
		&tweets,
		`
			SELECT id, message, tag, created_at
			FROM Tweets
			WHERE tag = ?
			LIMIT ? OFFSET ?
		`,
		tag, limit, offset,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get tweets: %w", err)
	}

	return tweets, nil
}

func (t TwitterDatabase) AggregateTweetsByYear(ctx context.Context, from time.Time, to time.Time) ([]models.YearlyAggregate, error) {
	var aggregates []models.YearlyAggregate
	err := t.db.SelectContext(
		ctx,
		&aggregates,
		`
			SELECT YEAR(created_at) as year, count(id) as tweets
			FROM Tweets
			WHERE created_at BETWEEN ? AND ?
			GROUP BY year
			ORDER BY year ASC
		`,
		from, to,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to aggregate tweets: %w", err)
	}

	return aggregates, nil
}

func (t TwitterDatabase) AggregateTweetsByMonth(ctx context.Context, from time.Time, to time.Time) ([]models.MonthlyAggregate, error) {
	var aggregates []models.MonthlyAggregate
	err := t.db.SelectContext(
		ctx,
		&aggregates,
		`
			SELECT YEAR(created_at) as year, MONTH(created_at) as month, count(id) as tweets
			FROM Tweets
			WHERE created_at BETWEEN ? AND ?
			GROUP BY year, month
			ORDER BY year, month ASC
		`,
		from, to,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to aggregate tweets: %w", err)
	}

	return aggregates, nil
}

func NewTwitterDatabase(db DB) TwitterDatabase {
	return TwitterDatabase{db: db}
}
