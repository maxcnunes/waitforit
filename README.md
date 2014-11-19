# wait for it

Wait for a port be available in a specific host.


### Options

- **`-host`**: Host to connect
- **`-port`**: Port to connect
- **`-timeout`**: Timeout to wait port be available


### Example

```bash
waitforit -host=google.com -port=90 -timeout=20
```


## Development

Running with `Docker` and `Fig`:

```bash
fig run --rm local
```
