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

```bash
waitforit -host=google.com -port=90 -timeout=20 -debug

waitforit -full-connection=tcp://google.com:90 -timeout=20 -debug

waitforit -full-connection=http://google.com -timeout=20 -debug

waitforit -full-connection=http://google.com:90 -timeout=20 -debug
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
