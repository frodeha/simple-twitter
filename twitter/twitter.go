package twitter

import (
	"context"
	"simple_twitter/models"
	"unicode/utf8"
)

const (
	MAX_PAGE_SIZE                 = 500 // Max number of tweets you can request at once
	MAX_TWEET_MESSAGE_LENGTH_UTF8 = 160 // Max length of a tweet message (UTF8 length)
	MAX_TWEET_TAG_LENGTH          = 32  // Max length of tag (byte length)
)

type Twitter struct {
	tweets TweetStorage
}

type TweetStorage interface {
	GetTweet(ctx context.Context, id int64) (models.Tweet, error)
	ListTweets(ctx context.Context, tag string, offset int, limit int) ([]models.Tweet, error)
	CreateTweet(ctx context.Context, message string, tag string) (int64, error)
}

func (t Twitter) CreateTweet(ctx context.Context, message string, tag string) (models.Tweet, error) {
	err := validateMessage(message)
	if err != nil {
		return models.Tweet{}, err
	}

	err = validateTag(tag)
	if err != nil {
		return models.Tweet{}, err
	}

	id, err := t.tweets.CreateTweet(ctx, message, tag)
	if err != nil {
		return models.Tweet{}, models.ErrInternalWithCause("failed to create tweet", err)
	}

	tweet, err := t.tweets.GetTweet(ctx, id)
	if err != nil {
		return models.Tweet{}, models.ErrInternalWithCause("failed to create tweet", err)
	}

	return tweet, nil
}

func (t Twitter) ListTweets(ctx context.Context, tag string, offset int, limit int) ([]models.Tweet, error) {
	if offset < 0 {
		return nil, models.ErrInvalid("offset can't be negative")
	}

	if limit > MAX_PAGE_SIZE {
		limit = MAX_PAGE_SIZE
	}

	tweets, err := t.tweets.ListTweets(ctx, tag, offset, limit)
	if err != nil {
		return nil, models.ErrInternalWithCause("failed to create tweet", err)
	}

	return tweets, nil
}

func validateTag(tag string) error {
	if tag == "" {
		return models.ErrInvalid("tag can't be empty")
	}

	if len(tag) > MAX_TWEET_TAG_LENGTH {
		return models.ErrInvalidf("tag is too long: max byte length is %d", MAX_TWEET_TAG_LENGTH)
	}

	return nil
}

func validateMessage(message string) error {
	if message == "" {
		return models.ErrInvalid("message can't be empty")
	}

	if utf8.RuneCountInString(message) > MAX_TWEET_MESSAGE_LENGTH_UTF8 {
		return models.ErrInvalidf("message is too long: max length is %d", MAX_TWEET_MESSAGE_LENGTH_UTF8)
	}

	return nil
}

func NewTwitter(tweets TweetStorage) Twitter {
	return Twitter{tweets}
}
