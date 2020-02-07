# checker
Basic HTTP Ping with modulable alerting

Checker will alert you if an URL you're monitoring returns an HTTP response status code that is not 200.

## Installation

Easiest way to use this is to compile it in a docker container :
```
docker build -t checker:latest -f ./Dockerfile ./
```

Populate .env file from .env.dist template. `MONITORED_URLS` should be a comma separated, http/https prefixed, set of URLs.

## Running checker

Then, running the script itself is a matter of populating environment variables in `.env` and starting docker container with environment file mounted

```
docker run --rm -v $(pwd)/.env:/go/bin/.env checker:latest --notifier=slack check
```

The official image built from this repository is also available on Docker hub image registry [https://cloud.docker.com/u/pauulog/repository/docker/pauulog/checker](https://cloud.docker.com/u/pauulog/repository/docker/pauulog/checker)

### Cron usage

In the following example, the program will check every minute if the URLs are up.
```
* * * * docker run --rm -v $(pwd)/.env:/go/bin/.env docker.io/pauulog/checker:latest --notifier=slack check
```
