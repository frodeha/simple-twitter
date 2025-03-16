# Preparation for production - Simple "Twitter" App 

## Background

During the production of this task I took left out some important features that is necessary before it can be deployed to production:

### Authentication

The "Aggregation endpoint" is supposed to be a "management endpoint" which I choose to interpret as an endpoint that shouldn't be publicly accessible. We would in that case need to implement some sort of authentication requirement on this endpoint so we don't expose this information to unintended audiences. I considered this to be outside the scope of this task, but would most certainly have to be implemented before the service can be deployed to production.

### Logging

There is no well structured and thought through system for logging in this server. A lot of the building blocks are there for it (custom errors that capture the cause of an error for example), but an actual logging utility will be necessary. 

If I were to implement this before pushing to production (which I very much would think necessary) I would implement a middleware in the `api` package that, for each request, instantiates an instance of a structured logger (`logrus` is my goto) that is populated with relevant request scoped log information (requestID, endpoint, IP perhaps etc).

I would then update the code to use this request scoped logger to log, for example, the causes of errors.

### Metrics

Lastly, we would need to have some sort of way to measure the load on the servers (number of requests per time frame, response time, error rates etc). This is important get a view of how our server performs, our "hot paths" and where we might have performance challenges. 

A pretty standard way of doing this would be to add metrics counters again as middlewares in the `api` package, and expose endpoints for something like `prometheus` to scrape.


## Deploying to production

I'm going to assume that the production environment is capable of running dockerized workloads (e.g docker containers) and we should therefore dockerize our app. This is a relatively straight forward process of:

* Adding a Dockerfile
    * I prefer to do a multistage build where the first step builds the application binary using a golang builder image and a secondary step that is based on something minimal (scratch or alpine) that we copy the application binary into. 
    * The latter image is the build output and should be as small and locked down as possible (dedicated minimal access users, no dependencies etc)

* Establish a CI/CD pipeline for testing, building, and deploying our code
    * There's many systems that lets you do this. I have good experience with Gitlab CI/CD, but Github actions or any of many others will do the job as well.
    * I'd like to have sequential steps of:
        * End to end tests
        * Build docker image and upload to some container registry
        * Deploy to a test environment if available
        * Deploy to production if a "release" (git tag) build
    
In production we should run with (at a minimum) of two replicas of the service for a basic level of redundancy. If we are expecting high amounts of traffic more replicas is a good idea. We should keep close attention on or metrics to see how the system is performing.

My hypothesis is that the first thing that is going to start struggling to handle high load is our database. We are querying our database on every request and that won't scale as the number of users grow.

Working from this hypothesis we have several tools at our disposal:

* Place a caching cluster (redis for example) in between our application and the database. This can be easily achieved in code by implementing a redis backed version of the interface required by the `twitter` package that checks the redis cluster before querying the database. On a cache miss (e.g, we have to query the database for the data), we update the cache so subsequent requests for this data can be served from cache for a while. This would reduce the amount of read queries towards our database substantially.

* Configure our database in a "single writer, many reader" configuration where one database replica is the writer (e.g primary) and we have a number of reader replicas that all are replicated to continuously from the single primary. If we update our server code to serve read queries (that miss the cache) from the read replicas and perform our inserts (writes) on the primary, we should again reduce load on the primary by a substantial amount. We have dedicated implementations for reading and writing in the `database` package, so this should be fairly straight forward (although it would introduce some replication lag in the system).

* Lastly, we could look to sharding our data. In the current API structure, reads and writes (except the aggregation) are all in context of a "tag". This could, if we need to, become a relatively decent shard key. We could use this to replicate the setup above as many times as we need and keep messages with the same "tag" (or that shard to the same shard based on tag) together in a database cluster and apply all the same mechanisms as above in addition to this. With this I suspect we can scale the system to very high loads. This would be a more involved coding task, but I believe we could contain it to the `database` package with the rest of the server being none the wiser.

