package test

import (
	"net/http"
	"simple_twitter/models"
	"strings"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (e *E2ETestSuite) Test_CreateTweet() {
	var (
		require = require.New(e.T())
		assert  = assert.New(e.T())
	)

	input := models.Tweet{
		Message: "This is a test tweet! ✅",
		Tag:     "e2e-tests",
	}

	res, err := http.Post(e.buildURL("/tweets", nil), "application/json", e.marshalTweet(input))
	require.NoError(err)
	defer res.Body.Close()

	assert.Equal(http.StatusCreated, res.StatusCode)
	output := e.unmarshalTweet(res)

	assert.NotZero(output.ID, "Expected created tweet to have an `id`")
	assert.Greater(output.ID, int64(2000), "Expected `id` of created tweet to be higher than 2000") // Seed data uses the first 2000 ID's
	assert.Equal(input.Message, output.Message, "Expected `message` to be the same in input and output")
	assert.Equal(input.Tag, output.Tag, "Expected `tag` to be the same in input and output")
	assert.NotZero(output.CreatedAt, "Expected `created at` of created tweet to be set")
}

func (e *E2ETestSuite) Test_CreateTweetWithoutMessage() {
	var (
		require = require.New(e.T())
		assert  = assert.New(e.T())
	)

	input := models.Tweet{Tag: "e2e-tests"}
	res, err := http.Post(e.buildURL("/tweets", nil), "application/json", e.marshalTweet(input))
	require.NoError(err)
	defer res.Body.Close()

	assert.Equal(http.StatusBadRequest, res.StatusCode, "Expected `status code` to be `400`")
	output := e.unmarshalError(res)

	assert.Equal(models.ErrKindInvalid, output.Kind, "Expected `error kind` to be `invalid`")
	assert.True(strings.Contains(output.Message, "message"), "Expected `error message` to contain `message`")
}

func (e *E2ETestSuite) Test_CreateTweetWithTooLongMessage() {
	var (
		require = require.New(e.T())
		assert  = assert.New(e.T())
	)

	input := models.Tweet{
		Message: "👋 This is a test tweet that shouldn't have a message that's longer than 140 code points! Is this longer than that? 🤔 Yeah - looks that way! ✅",
		Tag:     "e2e-tests",
	}

	res, err := http.Post(e.buildURL("/tweets", nil), "application/json", e.marshalTweet(input))
	require.NoError(err)
	defer res.Body.Close()

	assert.Equal(http.StatusBadRequest, res.StatusCode, "Expected `status code` to be `400`")
	output := e.unmarshalError(res)

	assert.Equal(models.ErrKindInvalid, output.Kind, "Expected `error kind` to be `invalid`")
	assert.True(strings.Contains(output.Message, "message"), "Expected `error message` to contain `tag`")
	assert.True(strings.Contains(output.Message, "too long"), "Expected `error message` to contain `too long`")
	assert.True(strings.Contains(output.Message, "140"), "Expected `error message` to contain `140`")
}

func (e *E2ETestSuite) Test_CreateTweetWithoutTag() {
	var (
		require = require.New(e.T())
		assert  = assert.New(e.T())
	)

	input := models.Tweet{Message: "This is a test tweet! ✅"}
	res, err := http.Post(e.buildURL("/tweets", nil), "application/json", e.marshalTweet(input))
	require.NoError(err)
	defer res.Body.Close()

	assert.Equal(http.StatusBadRequest, res.StatusCode, "Expected `status code` to be `400`")
	output := e.unmarshalError(res)

	assert.Equal(models.ErrKindInvalid, output.Kind, "Expected `error kind` to be `invalid`")
	assert.True(strings.Contains(output.Message, "tag"), "Expected `error message` to contain `tag`")
}

func (e *E2ETestSuite) Test_CreateTweetWithTooLongTag() {
	var (
		require = require.New(e.T())
		assert  = assert.New(e.T())
	)

	input := models.Tweet{
		Message: "This is a test tweet! ✅",
		Tag:     "e2e-tests-should-test-that-tags-cant-be-this-long",
	}

	res, err := http.Post(e.buildURL("/tweets", nil), "application/json", e.marshalTweet(input))
	require.NoError(err)
	defer res.Body.Close()

	assert.Equal(http.StatusBadRequest, res.StatusCode, "Expected `status code` to be `400`")
	output := e.unmarshalError(res)

	assert.Equal(models.ErrKindInvalid, output.Kind, "Expected `error kind` to be `invalid`")
	assert.True(strings.Contains(output.Message, "tag"), "Expected `error message` to contain `tag`")
	assert.True(strings.Contains(output.Message, "too long"), "Expected `error message` to contain `too long`")
	assert.True(strings.Contains(output.Message, "32"), "Expected `error message` to contain `32`")
}
