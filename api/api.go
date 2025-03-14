package api

import (
	"context"
	"encoding/json"
	"net/http"
	"simple_twitter/models"
	"time"

	"strconv"
)

type TwitterService interface {
	CreateTweet(ctx context.Context, message string, tag string) (models.Tweet, error)
	ListTweets(ctx context.Context, tag string, offset int, limit int) ([]models.Tweet, error)
	AggregateTweets(ctx context.Context, from time.Time, to time.Time, groupBy string) (models.AggregatedTweets, error)
}

func NewServer(addr string, twitter TwitterService) http.Server {
	var mux http.ServeMux
	mux.HandleFunc("POST /tweets", createTweet(twitter))
	mux.HandleFunc("GET /tweets", listTweets(twitter))
	mux.HandleFunc("GET /tweets/_aggregate", aggregateTweets(twitter))
	return http.Server{
		Addr:    addr,
		Handler: &mux,
	}
}

func createTweet(twitter TwitterService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var t models.Tweet
		err := json.NewDecoder(r.Body).Decode(&t)
		if err != nil {
			handleError(models.ErrInvalidWithCause("failed to parse request body", err), w, r)
			return
		}

		tweet, err := twitter.CreateTweet(r.Context(), t.Message, t.Tag)
		if err != nil {
			handleError(err, w, r)
			return
		}

		writeJSONResponse(http.StatusCreated, tweet, w)
	}
}

func listTweets(twitter TwitterService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			tag    = r.URL.Query().Get("tag")
			offset = 0
			limit  = 50
		)

		if r.URL.Query().Has("offset") {
			o, err := strconv.Atoi(r.URL.Query().Get("offset"))
			if err != nil {
				handleError(models.ErrInvalidWithCause("`offset` must be an integer value", err), w, r)
				return
			}
			offset = o
		}

		if r.URL.Query().Has("limit") {
			l, err := strconv.Atoi(r.URL.Query().Get("limit"))
			if err != nil {
				handleError(models.ErrInvalidWithCause("`limit` must be an integer value", err), w, r)
				return
			}
			limit = l
		}

		tweets, err := twitter.ListTweets(r.Context(), tag, offset, limit)
		if err != nil {
			handleError(err, w, r)
			return
		}

		writeJSONResponse(http.StatusOK, tweets, w)
	}
}

func aggregateTweets(twitter TwitterService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			groupBy = r.URL.Query().Get("group_by")
			from    time.Time
			to      time.Time
		)

		if r.URL.Query().Has("from") {
			f, err := time.Parse(time.DateOnly, r.URL.Query().Get("from"))
			if err != nil {
				handleError(models.ErrInvalidWithCause("`from` must be a valid date (YYYY-MM-DD)", err), w, r)
				return
			}
			from = f
		}

		if r.URL.Query().Has("to") {
			t, err := time.Parse(time.DateOnly, r.URL.Query().Get("to"))
			if err != nil {
				handleError(models.ErrInvalidWithCause("`to` must be a valid date (YYYY-MM-DD)", err), w, r)
				return
			}
			to = t
		}

		tweets, err := twitter.AggregateTweets(r.Context(), from, to, groupBy)
		if err != nil {
			handleError(err, w, r)
			return
		}

		writeJSONResponse(http.StatusOK, tweets, w)
	}
}

func handleError(err error, w http.ResponseWriter, r *http.Request) {
	statusCode := http.StatusInternalServerError
	if e, ok := err.(models.Error); ok {
		switch e.Kind {
		case models.ErrKindMissing:
			statusCode = http.StatusNotFound
		case models.ErrKindInvalid:
			statusCode = http.StatusBadRequest
		}
	}

	writeJSONResponse(statusCode, err, w)
}

func writeJSONResponse(statusCode int, response interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		// Todo(frode): Nothing much to do but log here
	}
}
