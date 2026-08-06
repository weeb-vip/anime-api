# anime-api

The anime catalogue service for weeb.vip: the series themselves, and everything
the site shows about them.

A Go GraphQL service built with [gqlgen](https://gqlgen.com) with federation
support, so it composes into the gateway schema alongside the other services
rather than being called directly. MySQL underneath, through GORM.

## Running it

Requires Go and MySQL.

```sh
make migrate                  # bring the database up to date
go run cmd/main.go server     # the GraphQL server
```

`config/config.dev.json` is the local config; `config/config.go` defines the
shape, and environment variables override it.

## Schema and generated code

```sh
make gql        # regenerate resolvers from graph/schema.graphqls
make mocks      # regenerate test mocks
make generate   # both
```

Generated output is committed, so regenerate and commit it with any schema
change.

## Migrations

```sh
make migrate-create name=add_something
make migrate
```

Migrations live in `db/migrations` as up/down SQL pairs.

## Metrics

Instrumented through [go-metrics-lib](https://github.com/weeb-vip/go-metrics-lib),
which reports to Datadog and Prometheus.

## Notes

Several rounds of query and index work are written up in
`DATABASE_OPTIMIZATION.md`, `OPTIMIZATION_NOTES.md` and
`PERFORMANCE_OPTIMIZATION_SUMMARY.md`; `ENVIRONMENT_LOGGING.md` covers how
logging differs between environments.

`.mirrord/` is configured for [mirrord](https://mirrord.dev), so a local
process can run against the cluster's traffic and dependencies.
