package test

import (
	"net/http"
	"net/url"
	"simple_twitter/models"
	"strings"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (e *E2ETestSuite) Test_AggregateTweetsByYear() {
	var (
		require = require.New(e.T())
		assert  = assert.New(e.T())
	)

	res, err := http.Get(e.buildURL("/tweets/_aggregate", url.Values{"group_by": {"year"}, "from": {"2024-01-01"}, "to": {"2025-12-31"}}))
	require.NoError(err)
	defer res.Body.Close()

	assert.Equal(http.StatusOK, res.StatusCode)
	aggregate := e.unmarshalAggregate(res)

	assert.Equal("year", aggregate.GroupBy, "Expected `group by` to be `year`")
	for idx, aggregate := range aggregate.Aggregates {
		assert.IsType(models.YearlyAggregate{}, aggregate, "Expected `aggregate` to be of yearly type")

		yearlyAggregate := aggregate.(models.YearlyAggregate)
		assert.Equal(2024+idx, yearlyAggregate.Year, "Expected `year` to be in ascending order")

		if yearlyAggregate.Year == 2024 {
			assert.Equal(770, yearlyAggregate.Tweets, "Expected `tweets` for 2024 to be `770`")
		}

		if yearlyAggregate.Year == 2025 {
			assert.GreaterOrEqual(yearlyAggregate.Tweets, 1022, "Expected `tweets` for 2025 to be at least `1022`")
		}
	}
}

func (e *E2ETestSuite) Test_AggregateTweetsByMonth() {
	var (
		require = require.New(e.T())
		assert  = assert.New(e.T())
	)

	res, err := http.Get(e.buildURL("/tweets/_aggregate", url.Values{"group_by": {"month"}, "from": {"2024-01-01"}, "to": {"2024-12-31"}}))
	require.NoError(err)
	defer res.Body.Close()

	assert.Equal(http.StatusOK, res.StatusCode)
	aggregates := e.unmarshalAggregate(res)

	assert.Equal("month", aggregates.GroupBy, "Expected `group by` to be `month`")
	for idx, aggregate := range aggregates.Aggregates {
		assert.IsType(models.MonthlyAggregate{}, aggregate, "Expected `aggregate` to be of monthly type")

		monthlyAggregate := aggregate.(models.MonthlyAggregate)
		assert.Equal(2024, monthlyAggregate.Year, "Expected `year` to be 2024")
		if idx > 0 {
			assert.Greater(
				monthlyAggregate.Month,
				aggregates.Aggregates[idx-1].(models.MonthlyAggregate).Month,
				"Expected `month` to be in ascending order",
			)
		}

		var expected int
		switch monthlyAggregate.Month {
		case 3:
			expected = 40
		case 4:
			expected = 69
		case 5:
			expected = 87
		case 6:
			expected = 58
		case 7:
			expected = 86
		case 8:
			expected = 99
		case 9:
			expected = 80
		case 10:
			expected = 91
		case 11:
			expected = 66
		case 12:
			expected = 92
		}

		assert.Equalf(expected, monthlyAggregate.Tweets, "Expected `tweets` for month `%d` to be `%d`", monthlyAggregate.Month, expected)
	}
}

func (e *E2ETestSuite) Test_AggregateTweetsByInvalidFromAndTo() {
	var (
		require = require.New(e.T())
		assert  = assert.New(e.T())
	)

	res, err := http.Get(e.buildURL("/tweets/_aggregate", url.Values{"group_by": {"month"}, "from": {"2025-01-01"}, "to": {"2024-01-01"}}))
	require.NoError(err)
	defer res.Body.Close()

	assert.Equal(http.StatusBadRequest, res.StatusCode)
	output := e.unmarshalError(res)

	assert.Equal(models.ErrKindInvalid, output.Kind, "Expected `error kind` to be `invalid`")
	assert.True(strings.Contains(output.Message, "from"), "Expected `error message` to contain `from`")
	assert.True(strings.Contains(output.Message, "to"), "Expected `error message` to contain `to`")
}

func (e *E2ETestSuite) Test_AggregateTweetsByInvalidGroup() {
	var (
		require = require.New(e.T())
		assert  = assert.New(e.T())
	)

	res, err := http.Get(e.buildURL("/tweets/_aggregate", url.Values{"group_by": {"day"}, "from": {"2024-01-01"}, "to": {"2024-12-31"}}))
	require.NoError(err)
	defer res.Body.Close()

	assert.Equal(http.StatusBadRequest, res.StatusCode)
	output := e.unmarshalError(res)

	assert.Equal(models.ErrKindInvalid, output.Kind, "Expected `error kind` to be `invalid`")
	assert.True(strings.Contains(output.Message, "group by"), "Expected `error message` to contain `group by`")
}
