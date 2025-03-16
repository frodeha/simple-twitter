# Simple "Twitter" App

This is a simple "Twitter flavoured" API server implementation in Go, done as a step in a technical interview process.
The task was to implement a server that allows clients to:

### Post messages with a tag
```bash
POST /tweets { "message": "This is a very interesting tweet 👍", "tag": "interesting-stuff" }

{
    "id": 2001,
    "message": "This is a very interesting tweet 👍",
    "tag": "interesting-stuff"
    "created_at": "2025-03-16T18:13:11Z"
}
```

### List messages with a given tag
```bash
GET /tweets?tag=interesting-stuff&offset=0&limit=50
[
    {
        "id": 2001,
        "message": "This is a very interesting tweet 👍",
        "tag": "interesting-stuff"
        "created_at": "2025-03-16T18:13:11Z"
    }
]
```

### Aggregate and count tweets posted in a given time period
```bash
GET /tweets/_aggregate?group_by=year&from=2024-01-01&to=2025-12-31
{
  "group_by": "year",
  "aggregates": [
    {
      "year": 2024,
      "tweets": 770
    },
    {
      "year": 2025,
      "tweets": 593
    }
  ]
}
```

The code is structured into packages according to a reasonable "division of responsibilities" mindset. The three main packages are `api` (responsible for the HTTP api), `twitter` (responsible for the business logic) and `database` (responsible for the data storage and retrieval). Packages define the interfaces they expect to receive in their respective constructors and implementations are instantiated and injected in `cmd/server/main.go`.

The `models` package holds the shared definitions of the domain types and the respective packages use these types in their interfaces. This way the packages can communicate using shared types without knowing anything about each other resulting in a loosely coupled codebase.

## Running

You will need a couple of tools to run this system effectively:
* The Go programming language (version 1.24 or higher)
* Docker and docker compose
* Make (you can run the commands in `Makefile` manually if you don't have this)

The server runs required dependencies (database) in docker compose. When the compose stack is brought up, it will 
  1) create a database connected to localhost port 3308, 
  2) apply database schema migration files and 
  3) seed the database with test data. 

Lastly it will run the server code on your host computer while using the database running in the docker compose stack. This approach is the expected way of running both the code and the tests (see below).

Simply issuing `make run` should set everything up for you:
```bash
$ make run
go build -o build/simple-twitter cmd/server/main.go
docker compose up -d
[+] Running 4/4
 ✔ Network simple-twitter_default      Created                                   0.0s
 ✔ Container simple-twitter-mysql-1    Healthy                                  30.7s
 ✔ Container simple-twitter-migrate-1  Exited                                   31.3s
 ✔ Container simple-twitter-seed-1     Started                                  31.4s
./build/simple-twitter
```

The server should now be ready to accept incoming requests on `localhost:3000`. Some handy `curl` commands that can be copy/paste'ed:

```bash
# Create a tweet
curl "localhost:3000/tweets" -XPOST -d '{ "message":"Hello world!", "tag":"greetings" }'

# List tweets
curl "localhost:3000/tweets?tag=greetings&offset=0&limit=50"

# Aggregate tweets by year
curl -s "localhost:3000/tweets/_aggregate?from=2022-01-01&to=2025-07-31&group_by=year"

# Aggregate tweets by month
curl -s "localhost:3000/tweets/_aggregate?from=2025-01-01&to=2025-07-31&group_by=month"
```

When you're done running the server you can take down the docker compose stack by running:

```bash
$ docker compose down
[+] Running 4/3
 ✔ Container simple-twitter-seed-1     Removed                                  0.0s 
 ✔ Container simple-twitter-migrate-1  Removed                                  0.0s 
 ✔ Container simple-twitter-mysql-1    Removed                                  3.6s 
 ✔ Network simple-twitter_default      Removed                                  0.1s 
```

## Testing

The codebase comes with a relatively exhaustive set of end to end tests that exercise the entire codebase (api + business logic + database). These also run with the expectation that docker compose will be able to spin up and maintain a ready to go database and can be executed in a similar way as [above](#running).

```bash
$ make e2e-tests
docker compose up -d
[+] Running 3/3
 ✔ Container simple-twitter-mysql-1    Healthy                                  0.5s
 ✔ Container simple-twitter-migrate-1  Exited                                   1.1s
 ✔ Container simple-twitter-seed-1     Started                                  1.3s
go test ./... -count=1
?       simple_twitter/api      [no test files]
?       simple_twitter/cmd/seed-script  [no test files]
?       simple_twitter/cmd/server       [no test files]
?       simple_twitter/database [no test files]
ok      simple_twitter/e2e_test 0.294s
?       simple_twitter/models   [no test files]
?       simple_twitter/twitter  [no test files]
```