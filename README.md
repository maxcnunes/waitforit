# wait for it

Wait for a port be available in a specific host.

![](http://24.media.tumblr.com/tumblr_m3x648wxbj1ru99qvo1_500.png)


### Options

- **`-host`**: Host to connect
- **`-port`**: Port to connect
- **`-timeout`**: Timeout to wait port be available
- **`-debug`**: Enable debug


### Example

```bash
waitforit -host=google.com -port=90 -timeout=20 -debug
```


## Development

Running with `Docker` and `Fig`:

```bash
fig run --rm local
```
