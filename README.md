# wait for it

[![Build Status](https://travis-ci.org/maxcnunes/waitforit.svg?branch=master)](https://travis-ci.org/maxcnunes/waitforit)
[![Coverage Status](https://coveralls.io/repos/github/maxcnunes/waitforit/badge.svg?branch=master)](https://coveralls.io/github/maxcnunes/waitforit?branch=master)

Wait until an address become available.

![](http://24.media.tumblr.com/tumblr_m3x648wxbj1ru99qvo1_500.png)


### Download

[Releases](https://github.com/maxcnunes/waitforit/releases)

### Options

- `-address`: Address (e.g. http://google.com or tcp://mysql_ip:mysql_port) - *former **full-connection***
- `-host`: Host to connect
- `-port`: Port to connect (default 80)
- `-timeout`: Time to wait until the address become available
- `-debug`: Enable debug
- `-v`: Show the current version
- `-file`: Path to the JSON file with the configs
- `-- `: Execute a post command once the address became available

### Example

#### Running

```bash
waitforit -host=google.com -port=90 -timeout=20 -debug

waitforit -address=tcp://google.com:90 -timeout=20 -debug

waitforit -address=http://google.com -timeout=20 -debug

waitforit -address=http://google.com:90 -timeout=20 -debug

waitforit -address=http://google.com -timeout=20 -debug -- printf "Google Works\!"
```

#### Using with config file

Create a JSON file describing the hosts you would like to wait for.

Example JSON:
```json
{
  "configs": [
    {
      "host": "google.com",
      "port": 80,
      "timeout": 20
    },
    {
      "address": "http://google.com:80",
      "timeout": 40
    }
  ]
}
```

```bash
waitforit -file=./config.json
```

#### Installing with a Dockerfile

##### Using curl

```
FROM node:6.5.0

ENV WAITFORIT_VERSION="v1.3.1"
RUN curl -o /usr/local/bin/waitforit -sSL https://github.com/maxcnunes/waitforit/releases/download/$WAITFORIT_VERSION/waitforit-linux_amd64 && \
    chmod +x /usr/local/bin/waitforit
```

##### Using wget

```
FROM node:6.5.0

ENV WAITFORIT_VERSION="v1.3.1"
RUN wget -q -O /usr/local/bin/waitforit https://github.com/maxcnunes/waitforit/releases/download/$WAITFORIT_VERSION/waitforit-linux_amd64 \
    && chmod +x /usr/local/bin/waitforit
```

##### Using COPY (from local file system)

```
FROM node:6.5.0

COPY waitforit-linux_amd64 /usr/local/bin/waitforit
RUN chmod +x /usr/local/bin/waitforit
```

## Development

```bash
make run
```

Running with `Docker` and `Compose`:

```bash
docker-compose run --rm local
```

```bash
docker-compose run --rm local go run src/waitforit/main.go -h
```

## Test

```bash
make test
```

## Build

```bash
make build
```
