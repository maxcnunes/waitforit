# wait for it

Wait until a port become available in a specific host.

![](http://24.media.tumblr.com/tumblr_m3x648wxbj1ru99qvo1_500.png)


### Download

[Releases](https://github.com/maxcnunes/waitforit/releases)

### Options

- **-full-connection**: Full connection `<protocol>://<host>:<port>`
- **-host**: Host to connect
- **-port**: Port to connect
- **-timeout**: Time to wait until port become available
- **-debug**: Enable debug


### Example

```bash
waitforit -host=google.com -port=90 -timeout=20 -debug
```

```bash
waitforit -full-connection=tcp://google.com:90 -timeout=20 -debug
```

## Development

Running with `Docker` and `Compose`:

```bash
docker-compose run --rm local
```
