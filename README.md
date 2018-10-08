# go-sample-provider

This repo contains a minimal provider using only Go. It builds bears as a service!

## Requirements

[Go](https://golang.org/) v 1.11+

## Setup

```
make install
```

## Testing

In one terminal start the server with:

```
go run cmd/server/main.go
```

And in another terminal run:

```
make test
```