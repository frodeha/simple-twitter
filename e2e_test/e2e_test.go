package test

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"simple_twitter/api"
	"simple_twitter/database"
	"simple_twitter/models"
	"simple_twitter/twitter"
	"testing"
	"time"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, &E2ETestSuite{})
}

type E2ETestSuite struct {
	suite.Suite

	conn   *sqlx.DB
	dbName string

	server *httptest.Server
}

func (e *E2ETestSuite) SetupSuite() {
	var (
		username = "root"
		password = "TopSecret"
		host     = "localhost:3308"

		require = require.New(e.T())
	)

	conn, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/?parseTime=true&multiStatements=true", username, password, host))
	require.NoError(err)

	err = conn.Ping()
	require.NoError(err)

	b := make([]byte, 8)
	rand.New(rand.NewSource(time.Now().Unix())).Read(b)
	dbName := hex.EncodeToString(b)

	// NB: Don't do this in production
	_, err = conn.Exec(fmt.Sprintf("CREATE DATABASE `%s`;", dbName))
	require.NoError(err)

	err = conn.Close()
	require.NoError(err)

	conn, err = sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&multiStatements=true", username, password, host, dbName))
	require.NoError(err)

	e.conn = conn
	e.dbName = dbName

	m, err := mysql.WithInstance(conn.DB, &mysql.Config{})
	require.NoError(err)

	migrations, err := migrate.NewWithDatabaseInstance("file://../database/migrations", "mysql", m)
	require.NoError(err)

	err = migrations.Up()
	require.NoError(err)

	m, err = mysql.WithInstance(conn.DB, &mysql.Config{MigrationsTable: "seed_migrations"})
	require.NoError(err)

	seeds, err := migrate.NewWithDatabaseInstance("file://../database/seeds", "mysql", m)
	require.NoError(err)

	err = seeds.Up()
	require.NoError(err)

	twitter := twitter.NewTwitter(database.NewTwitterDatabase(e.conn))
	server := api.NewServer("", twitter)

	e.server = httptest.NewServer(server.Handler)
}

func (e *E2ETestSuite) TearDownSuite() {
	var (
		require = require.New(e.T())
	)

	e.server.Close()

	// NB: Don't do this either in production
	_, err := e.conn.Exec(fmt.Sprintf("DROP DATABASE `%s`", e.dbName))
	require.NoError(err)

	err = e.conn.Close()
	require.NoError(err)
}

func (e *E2ETestSuite) buildURL(path string, query url.Values) string {
	url, err := url.Parse(e.server.URL)
	require.NoError(e.T(), err)

	url.Path = path
	url.RawQuery = query.Encode()

	return url.String()
}

func (e *E2ETestSuite) unmarshalTweet(res *http.Response) models.Tweet {
	var tweet models.Tweet
	err := json.NewDecoder(res.Body).Decode(&tweet)
	require.NoError(e.T(), err)
	return tweet
}

func (e *E2ETestSuite) marshalTweet(tweet models.Tweet) io.Reader {
	b, err := json.Marshal(tweet)
	require.NoError(e.T(), err)
	return bytes.NewReader(b)
}

func (e *E2ETestSuite) unmarshalTweets(res *http.Response) []models.Tweet {
	var tweets []models.Tweet
	err := json.NewDecoder(res.Body).Decode(&tweets)
	require.NoError(e.T(), err)
	return tweets
}

func (e *E2ETestSuite) unmarshalAggregate(res *http.Response) models.AggregatedTweets {
	var aggregate models.AggregatedTweets
	err := json.NewDecoder(res.Body).Decode(&aggregate)
	require.NoError(e.T(), err)
	return aggregate
}

func (e *E2ETestSuite) unmarshalError(res *http.Response) models.Error {
	var error models.Error
	err := json.NewDecoder(res.Body).Decode(&error)
	require.NoError(e.T(), err)
	return error
}
