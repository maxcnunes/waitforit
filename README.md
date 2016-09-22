# wait for it

Wait until an address become available.

![](http://24.media.tumblr.com/tumblr_m3x648wxbj1ru99qvo1_500.png)


### Download

[Releases](https://github.com/maxcnunes/waitforit/releases)

### Options

- **-full-connection**: Full connection `<protocol/scheme>://<host>:<port>`
- **-host**: Host to connect
- **-port**: Port to connect (default 80)
- **-timeout**: Time to wait until port become available
- **-debug**: Enable debug
- **-v**: Show the current version


### Example

#### Running

```bash
waitforit -host=google.com -port=90 -timeout=20 -debug

waitforit -full-connection=tcp://google.com:90 -timeout=20 -debug

waitforit -full-connection=http://google.com -timeout=20 -debug

waitforit -full-connection=http://google.com:90 -timeout=20 -debug
```

#### Installing with a Dockerfile

```
FROM node:6.5.0

ENV WAITFORIT_VERSION="v1.3.1"
RUN wget -q -O /usr/local/bin/waitforit https://github.com/maxcnunes/waitforit/releases/download/$WAITFORIT_VERSION/waitforit-linux_amd64 \
    && chmod +x /usr/local/bin/waitforit
```

## Development

Running with `Docker` and `Compose`:

```bash
docker-compose run --rm local
```

```bash
docker-compose run --rm local go run src/waitforit/main.go -h
```


## Build

Using [goxc](https://github.com/laher/goxc).

```bash
goxc
```
