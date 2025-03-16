package test

import (
	"net/http"
	"net/url"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (e *E2ETestSuite) Test_GetTweetsWithTag() {
	var (
		require = require.New(e.T())
		assert  = assert.New(e.T())
	)

	res, err := http.Get(e.buildURL("/tweets", url.Values{"tag": {"protocol-reboot"}}))
	require.NoError(err)
	defer res.Body.Close()

	assert.Equal(http.StatusOK, res.StatusCode)
	tweets := e.unmarshalTweets(res)
	for _, tweet := range tweets {
		assert.Equal("protocol-reboot", tweet.Tag, "Expected al tweets to have tag `protocol-reboot`")
	}
}

func (e *E2ETestSuite) Test_GetTweetsWithLimit() {
	var (
		require = require.New(e.T())
		assert  = assert.New(e.T())
	)

	res, err := http.Get(e.buildURL("/tweets", url.Values{"tag": {"protocol-reboot"}, "limit": []string{"2"}}))
	require.NoError(err)
	defer res.Body.Close()

	assert.Equal(http.StatusOK, res.StatusCode)
	tweets := e.unmarshalTweets(res)
	assert.Len(tweets, 2, "Expected 2 tweets to be returned when limit is 2")
	for _, tweet := range tweets {
		assert.Equal("protocol-reboot", tweet.Tag, "Expected al tweets to have tag `protocol-reboot`")
	}
}

func (e *E2ETestSuite) Test_GetTweetsWithNoTag() {
	var (
		require = require.New(e.T())
		assert  = assert.New(e.T())
	)

	res, err := http.Get(e.buildURL("/tweets", nil))
	require.NoError(err)
	defer res.Body.Close()

	assert.Equal(http.StatusOK, res.StatusCode)
	tweets := e.unmarshalTweets(res)
	assert.Len(tweets, 0, "Expected no tweets to be returned when tag is unspecified")
}
