# Updated Metrics System

## What's Changed

The metrics system has been centralized with default tags and simplified interfaces.

### ✅ Before (Old Way)
```go
// Verbose, repetitive, error-prone
_ = metrics.NewMetricsInstance().ResolverMetric(
    float64(time.Since(startTime).Milliseconds()),
    metrics_lib.ResolverMetricLabels{
        Resolver: "AnimeBySeasons",
        Service:  "anime-api",
        Protocol: "graphql",
        Result:   metrics_lib.Success,
        Env:      metrics.GetCurrentEnv(),
    },
)
```

### ✅ After (New Way)
```go
// Clean, simple, consistent
metrics.GetAppMetrics().ResolverMetric(
    float64(time.Since(startTime).Milliseconds()),
    "AnimeBySeasons",
    metrics.Success,
)
```

## Key Benefits

### 🎯 **Default Tags**
- `service`: Automatically set from config (`anime-api`)
- `env`: Automatically set from config (`development`/`production`)
- `version`: Automatically set from config
- `protocol`: Automatically set to `graphql` for resolvers

### 🧹 **Cleaner Interface**
- Only specify what changes: resolver name, table, method, result
- No more repetitive service/env/protocol parameters
- Consistent across all components

### 📊 **Centralized Configuration**
- Single metrics instance with proper initialization
- Default buckets and labels configured once
- Removed DataDog dependencies (Prometheus only)

## Usage Examples

### Resolver Metrics
```go
// Error case
metrics.GetAppMetrics().ResolverMetric(
    float64(time.Since(startTime).Milliseconds()),
    "AnimeBySeasons",
    metrics.Error,
)

// Success case
metrics.GetAppMetrics().ResolverMetric(
    float64(time.Since(startTime).Milliseconds()),
    "AnimeBySeasons",
    metrics.Success,
)
```

### Database Metrics
```go
// Database operation
metrics.GetAppMetrics().DatabaseMetric(
    float64(time.Since(startTime).Milliseconds()),
    metrics.TableAnime,      // or "anime"
    metrics.MethodSelect,    // or "SELECT"
    metrics.Success,
)
```

### Repository Metrics
```go
// Repository operation (uses DatabaseMetric internally)
metrics.GetAppMetrics().RepositoryMetric(
    float64(time.Since(startTime).Milliseconds()),
    "FindBySeasonWithFieldSelection",
    "find",
    metrics.Success,
)
```

## Constants Available

### Results
- `metrics.Success`
- `metrics.Error`

### Database Methods
- `metrics.MethodSelect`
- `metrics.MethodInsert`
- `metrics.MethodUpdate`
- `metrics.MethodDelete`

### Table Names
- `metrics.TableAnime`
- `metrics.TableAnimeSeason`
- `metrics.TableEpisodes`
- `metrics.TableCharacters`
- `metrics.TableStaff`

## Automatic Tags

All metrics automatically include:
```
service="anime-api"
env="development"      # from config
version="x.x.x"        # from config
protocol="graphql"     # for resolvers
```

You only specify the variable parts!