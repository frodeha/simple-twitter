package api

import (
	"context"
	"encoding/json"
	"net/http"
	"simple_twitter/models"

	"strconv"
)

type TwitterService interface {
	CreateTweet(ctx context.Context, message string, tag string) (models.Tweet, error)
	ListTweets(ctx context.Context, tag string, offset int, limit int) ([]models.Tweet, error)
}

func NewServer(addr string, twitter TwitterService) http.Server {
	var mux http.ServeMux

	mux.HandleFunc("POST /tweets", createTweet(twitter))
	mux.HandleFunc("GET /tweets", listTweets(twitter))

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
			handleError(err, w, r)
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
				handleError(err, w, r)
				return
			}
			offset = o
		}

		if r.URL.Query().Has("limit") {
			l, err := strconv.Atoi(r.URL.Query().Get("limit"))
			if err != nil {
				handleError(err, w, r)
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

func handleError(err error, w http.ResponseWriter, r *http.Request) {
	var (
		statusCode int
		response   interface{}
	)

	switch err {
	// Todo(frode): Handle more error codes
	default:
		statusCode = http.StatusInternalServerError
		response = map[string]string{"message": "Internal server error"}
	}

	writeJSONResponse(statusCode, response, w)
}

func writeJSONResponse(statusCode int, response interface{}, w http.ResponseWriter) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		// Todo(frode): Nothing much to do but log here
	}
}
