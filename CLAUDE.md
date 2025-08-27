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
- Tests are located in `metrics_lib/` subdirectory

## Architecture

This is a GraphQL anime API built with Go, using:
- **gqlgen** for GraphQL server generation with federation support
- **GORM** with MySQL for database operations
- **Cobra CLI** for command structure
- **Custom metrics library** (`./metrics_lib`) with Datadog and Prometheus support

### Key Structure
- `cmd/main.go` - Application entry point
- `internal/commands/` - CLI commands (serve, migrate up/down)
- `graph/` - GraphQL schema, resolvers, and generated code
- `internal/db/repositories/` - Data access layer with entities and repositories
- `internal/services/` - Business logic layer
- `internal/resolvers/` - GraphQL resolver implementations
- `db/migrations/` - Database migration files
- `metrics_lib/` - Custom metrics library (separate Go module)

### GraphQL Schema
Main queries support anime search, episodes, characters/staff, and various ranking queries (newest, top rated, most popular). Schema files are in `graph/*.graphqls`.

### Database
Uses MySQL with GORM. Entities include anime, episodes, characters, staff, and relations. Migrations use golang-migrate with SQL files.

### Configuration
Configuration loaded via `internal/db/config.go` with development config in `config/config.dev.json`.