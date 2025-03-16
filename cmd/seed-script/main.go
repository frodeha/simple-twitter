package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

type tweet struct {
	message   string
	tag       string
	createdAt time.Time
}

func main() {
	b, err := os.ReadFile("cmd/seed-script/faker.csv")
	fatal(err)

	records, err := csv.NewReader(bytes.NewReader(b)).ReadAll()
	fatal(err)

	var tweets []tweet
	for _, record := range records {
		createdAt, err := time.Parse(time.RFC3339, record[2])
		fatal(err)
		tweets = append(tweets, tweet{message: record[0], tag: record[1], createdAt: createdAt})
	}
	sort.Slice(tweets, func(i, j int) bool { return tweets[i].createdAt.Before(tweets[j].createdAt) })

	var builder strings.Builder
	builder.WriteString("INSERT INTO `Tweets` (message, tag, created_at) VALUES")
	builder.WriteString("\n")
	for idx, tweet := range tweets {
		builder.WriteString(fmt.Sprintf(`("%s", "%s", "%s")`, tweet.message, tweet.tag, tweet.createdAt.Format(time.DateTime)))
		if idx != len(tweets)-1 {
			builder.WriteString(",")
		} else {
			builder.WriteString(";")
		}
		builder.WriteString("\n")
	}

	err = os.WriteFile("database/seeds/20250314164410_tweets.up.sql", []byte(builder.String()), 0777)
	fatal(err)
}

func fatal(err error) {
	if err != nil {
		panic(err)
	}
}
