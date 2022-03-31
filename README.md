# foaas-api


## Overview

This repository runs a server that expose a `/message` endpoint and forward request to [fuck off as a service](https://www.foaas.com/).
It's implements a rate limit on the http server.

## Requirements

The server is implemented in [Go 1.16](https://go.dev). 
To install Go, follow the [instructions](https://go.dev/doc/install).

## Getting Started

- To run the server, go to the root folder and execute:

```
make all
make run
```

If you have docker installed and want to run it in a containerised environment, run:

```
make docker-build
make docker-run
```

After setting up the server, it can be tested with the following `curl`. 
The request must contain a header called `UserId`

Example:

```
curl -H 'UserId: "123"' localhost:4000/message
```

- To test the code and see the coverage, go to the root folder and execute:

```
make test
make coverage
```

## Customising the server

There are many arguments to customize the server:

- log-level, by default it's debug. Ex: it can be info. 
- rate-limit-enable, by default it's true. It's enable the rate limit feature.
- rate-limit-count, by default it's 5. It's the number of requests allowed in the window time.
- rate-limit-window-in-milliseconds, by default it's 10000. It's the window time to evaluate the number of requests. 
- timeout-in-milliseconds, by default it's 10000. It's the timeout of the API call to `foaas-api`.

Example:

```
./foaas-api serve \
    --log-level=info \
    --rate-limit-enable=true \
    --rate-limit-count=5 \
    --rate-limit-window-in-milliseconds=10000 \
    --timeout-in-milliseconds=10000
```
