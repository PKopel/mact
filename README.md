# MA/CT

Mock API/Call Translation - proxy server for testing changes to http servers API.

## Usage

- `make build`/`make install` to build/install `mact` executable form source
- `mact --config=<PATH TO CONFIG YAML> --port=<PORT TO LISTEN ON>` to run proxy server

Example config file:

```yaml
services:
  - host: "http://localhost:8080"
    endpoints:
      - path: /ping
        verb: "GET"
        statusCodes: # list of status codes for which responses will be processed
          - 200
        changes:
          - type: add # suported types are add, modify, remove
            field: example.field # json field path
            value: {"test": value} # json value (ignored for 'remove' type)
```
