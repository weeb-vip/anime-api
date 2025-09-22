# Environment Logging Implementation

## Summary

Added environment context to all application logs to improve debugging and monitoring across different environments (development, staging, production).

## Changes Made

### 1. Logger Configuration Updated

**File: `internal/logger/logger.go`**
- Added `Environment` field to logger configuration
- Added `WithEnvironment()` option function
- Updated global logger initialization to include environment field

### 2. Application Initialization Updated

**File: `internal/commands/serve.go`**
- Updated logger initialization to include environment from config
- Environment is loaded from config and passed to logger

### 3. Automatic Environment Inclusion

All logs throughout the application now automatically include the environment field:

- **Global logs**: `logger.Get().Info().Msg("message")`
- **Context logs**: `logger.FromCtx(ctx).Info().Msg("message")`
- **Database logs**: SQL queries via traced logger
- **GraphQL logs**: All resolver and middleware logs
- **HTTP logs**: Server startup and request logs

## Environment Configuration

The environment is read from:
1. `ENV` environment variable
2. Config file `AppConfig.Env` field
3. Defaults to `"development"` if not set

## Example Log Output

### Before
```json
{
  "level": "info",
  "service": "anime-api",
  "version": "1.0.0",
  "time": "2025-09-21T23:41:21-04:00",
  "message": "GraphQL query executed"
}
```

### After
```json
{
  "level": "info",
  "service": "anime-api",
  "version": "1.0.0",
  "environment": "production",
  "time": "2025-09-21T23:41:21-04:00",
  "message": "GraphQL query executed"
}
```

## Usage Examples

### Setting Environment

```bash
# Development (default)
go run cmd/main.go serve

# Production
ENV=production go run cmd/main.go serve

# Staging
ENV=staging go run cmd/main.go serve
```

### Log Filtering by Environment

With logging aggregation tools like ELK, Splunk, or DataDog:

```bash
# Filter production logs
environment:"production"

# Filter development logs
environment:"development"

# Filter specific service in production
service:"anime-api" AND environment:"production"
```

## Benefits

1. **Environment Visibility**: Instantly know which environment generated each log
2. **Debugging**: Easier troubleshooting across environments
3. **Monitoring**: Better alerting and dashboard filtering
4. **Compliance**: Environment tracking for audit trails
5. **Performance Analysis**: Compare performance across environments

## Backward Compatibility

- Existing log parsing is unaffected (just adds a new field)
- No breaking changes to existing logging calls
- All existing logs automatically get environment context

## Testing

Comprehensive tests verify:
- ✅ Environment field included in global logger
- ✅ Environment field included in context logger
- ✅ Custom fields preserved alongside environment
- ✅ Configuration loading works correctly

## Integration with Existing Systems

The environment field works seamlessly with:
- **Tracing**: Environment included alongside trace/span IDs
- **Metrics**: Consistent environment labeling across logs and metrics
- **Database Logs**: SQL queries include environment context
- **GraphQL Logs**: All resolvers log with environment context