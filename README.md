# Sweldo REST API

## Prerequisites

Go version: `go version go1.26.0 darwin/arm64`

Install `goose` and `sqlc`

```zsh
$ brew install goose sqlc
```

## How to run locally

Run docker via the terminal or you can use docker desktop

```zsh
$ docker-compose up -d
```

Run the code

```zsh
$ go run cmd/*.go
```

## How to create migrations

Simply run the following command:

```zsh
$ goose -s description_of_migration sql
$ Created new file: 00001_description_of_migration.sql
```

## Generate sqlc types

Run the command:

```zsh
$ sqlc generate
```

This will populate the `adapters/postgresql/sqlc/` directory with `*.sql.go` files.
