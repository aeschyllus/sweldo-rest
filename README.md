# Sweldo REST API

## Prerequisites

Go version: `go version go1.26.0 darwin/arm64`

Install `goose` and `sqlc`

On MacOS

```zsh
$ brew install goose sqlc
```

On Arch Linux

```zsh
// for goose
$ go install github.com/pressly/goose/v3/cmd/goose@latest

// for sqlc
$ sudo pacman -S goose sqlc
```

## Generate sqlc types

Run the command:

```zsh
$ sqlc generate
```

This will populate the `adapters/postgresql/sqlc/` directory with `*.sql.go` files.

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
$ goose create -s description_of_migration sql
$ Created new file: 00001_description_of_migration.sql
```
