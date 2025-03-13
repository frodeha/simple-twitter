package main

import (
	"flag"
	"log"
	"os"
	"simple_twitter/api"
	"simple_twitter/database"
	"simple_twitter/twitter"

	ff "github.com/peterbourgon/ff/v3"
)

func main() {
	fs := flag.NewFlagSet("simple-twitter", flag.ExitOnError)

	var (
		mysqlAddr     = fs.String("mysql-addr", "127.0.0.1:3308", "")
		mysqlUser     = fs.String("mysql-user", "root", "")
		mysqlPassword = fs.String("mysql-password", "TopSecret", "")
		mysqlDatabase = fs.String("mysql-database", "simple_twitter", "")

		listenAddr = fs.String("listen-addr", "localhost:3000", "")
	)

	err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarNoPrefix())
	if err != nil {
		log.Fatal(err)
	}

	conn, err := database.Connect(*mysqlAddr, *mysqlUser, *mysqlPassword, *mysqlDatabase)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	var (
		tweetStorage = database.NewTwitterDatabase(conn)
		twitter      = twitter.NewTwitter(tweetStorage)
		apiServer    = api.NewServer(*listenAddr, twitter)
	)

	apiServer.ListenAndServe()

}
