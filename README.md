# API Wiremock

This set up a wiremock service and a reverse proxy service to intercept calls to a third party API and cache the response via wiremock.

# Usage

### Run:
```
make start
```

### Run locally (wiremock running in the background with docker):
```
make dev
```

### Remove all captured stub mappings from mocks/
```
make clean
```

# Configuration

| Environment Variable | Description | Required |
| --- | --- |
| `API-WIREMOCK` | Defines the URL of the wiremock service| Y |
| `API-TARGET` | Defines the URL of the 3rd party API| Y |
| `PROXY_URL` | Defines proxy url to call target API if any| N |


# License

Code herein is licensed under [the permissive MIT license](./LICENSE)