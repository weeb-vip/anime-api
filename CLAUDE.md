# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

### Development
- `go run cmd/main.go serve` - Start the GraphQL API server
- `make gql` - Regenerate GraphQL code from schema using gqlgen
- `make generate` - Generate all code (mocks + GraphQL)
- `make migrate-create name=<migration_name>` - Create new database migration

### Database
- `go run cmd/main.go migrate up` - Apply all pending database migrations
- `go run cmd/main.go migrate down` - Rollback one database migration
- `docker-compose up -d` - Start MySQL and Redis services locally

### Testing
- `go test ./...` - Run all tests
- Tests are located in the various package subdirectories

## Architecture

This is a GraphQL anime API built with Go, using:
- **gqlgen** for GraphQL server generation with federation support
- **GORM** with MySQL for database operations
- **Cobra CLI** for command structure
- **External metrics library** (`github.com/weeb-vip/go-metrics-lib`) with Datadog and Prometheus support

### Key Structure
- `cmd/main.go` - Application entry point
- `internal/commands/` - CLI commands (serve, migrate up/down)
- `graph/` - GraphQL schema, resolvers, and generated code
- `internal/db/repositories/` - Data access layer with entities and repositories
- `internal/services/` - Business logic layer
- `internal/resolvers/` - GraphQL resolver implementations
- `db/migrations/` - Database migration files

### GraphQL Schema
Main queries support anime search, episodes, characters/staff, and various ranking queries (newest, top rated, most popular). Schema files are in `graph/*.graphqls`.

### Database
Uses MySQL with GORM. Entities include anime, episodes, characters, staff, and relations. Migrations use golang-migrate with SQL files.

### Configuration
Configuration loaded via `internal/db/config.go` with development config in `config/config.dev.json`.

### Caching
- **Single repository pattern**: Repositories conditionally use caching based on configuration
- **Environment control**: Set `CACHE_ENABLED=false` to disable caching (default: `true`)
- **Cache TTL**: 30min for anime data, 15min for episodes, 1hr for season lists
- **Cache coordination**: Updates to anime automatically invalidate related episode cache
- **Graceful fallback**: When Redis is unavailable or disabled, repositories bypass cache entirely
- **Unified interface**: Same repository interface whether caching is enabled or disabled

## Important Notes

### Dependencies
The application uses the external metrics library `github.com/weeb-vip/go-metrics-lib` which provides Datadog and Prometheus metrics support.

### Code Generation
- GraphQL code is generated via gqlgen - modify `graph/*.graphqls` files and run `make gql`
- Mocks are generated and included in the `make generate` target
- Always regenerate after schema changes

### Database Migrations
- Uses golang-migrate with SQL migration files in `db/migrations/`
- Sequential naming convention enforced by migrate tool
- Apply with `go run cmd/main.go migrate up`, rollback with `migrate down`