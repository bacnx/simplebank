# SimpleBank

A simple back-end project using Golang, PostgreSQL.

## Technologies

- [Golang](https://go.dev)
- PostgreSQL
- Docker
- [sqlc](https://sqlc.dev)
- [migrate](https://github.com/golang-migrate/migrate) _(for migrate database)_
- [Paseto](https://paseto.io) _(create token for authentication)_

## Setup local development

### Setup database

- Create network `bank-network`:

```bash
make network
```

- Start postgres container:

```bash
make postgres
```

- Create `simple_bank` database:

```bash
make createdb
```

- Run db migration up all versions ([install migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate#installation)):

```bash
make migrateup
```

- Run db migration up 1 version:

```bash
make migrateup1
```

- Run db migration down all versions:

```bash
make migratedown
```

- Run db migration down 1 version:

```bash
make migratedown1
```

### How to generate code

- Generate SQL CRUD with sqcl:

```bash
make sqlc
```

- Generate DB mock with [gomock](https://pkg.go.dev/go.uber.org/mock/gomock):

```bash
make mock
```

- Create a new db migration:

```bash
make new_migration name=<migration_name>
```

### How to run

- Run server:

```bash
make server
```

- Run test:

```bash
make test
```
