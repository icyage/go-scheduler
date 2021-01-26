# go-scheduler

Probe Running service

## Requirements

go-scheduler requires pgsql and rabbitmq and redis to function.

```$xslt
docker run -d --hostname rabbit --name rabbit -p 15672:15672 -p 5672:5672 rabbitmq:3-management
```
```$xslt
docker run --name scheduler-postgres -e POSTGRES_PASSWORD=scheduler -e POSTGRES_USER=scheduler -e POSTGRES_DB=scheduler -p 5432:5432 -d postgres
```
```$xslt
docker run --name redis -p 6379:6379 -d redis
```

you can also use the same DB instance currently being used by the `scheduler backend`

## Getting started

This project requires Go to be installed. On OS X with Homebrew you can just run `brew install go`.

Running it then should be as simple as:

```console
$ make
$ ./bin/go-scheduler
```

### go-scheduler server types

go-scheduler currently has 3 server parts  that can all run individually and scale on their own.

HTTP Server `go-scheduler serve` handles the http REST api to manage and comunicate with the scheduler
Scheduler `go-scheduler scheduler` runs the scheduler and makes sure the different cron jobs run at their interval
Launcher `go-scheduler probe launcher` consumes the trigger queue from the scheduler and launches probe workers

the scheduler requires a minimal database in PGSQL to hold all the cron jobs this can be init by running `go-scheduler scheduler dbinit`


### Testing

``make test``# go-scheduler
