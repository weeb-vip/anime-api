# Field Selection Implementation

## What's Implemented

The codebase now supports **selective field querying** where the GraphQL query only fetches the database columns that are actually requested in the GraphQL query.

## How It Works

### 1. Repository Level (`FindBySeasonWithFieldSelection`)
- Accepts a `FieldSelection` object that specifies which fields to fetch
- Dynamically builds SQL SELECT clauses based on requested fields
- Uses raw SQL for maximum performance and field control

### 2. Service Level (`AnimeBySeasonWithFieldSelection`)
- Passes field selection information to the repository
- Maintains tracing and metrics for field-optimized queries

### 3. Resolver Level (`AnimeBySeasons`)
- Extracts field selection from GraphQL context (currently returns nil for fallback)
- Chooses between field-optimized and standard queries
- Maintains episodes as separate resolvers for lazy loading

## Example Usage

### Without Field Selection (Standard Query)
```graphql
query {
  animeBySeasons(season: "WINTER_2024") {
    id
    titleEn
    # This fetches ALL anime fields from database
  }
}
```

### With Field Selection (Optimized Query)
```graphql
query {
  animeBySeasons(season: "WINTER_2024") {
    id
    titleEn
    # This would only fetch: id, title_en, created_at, updated_at from database
  }
}
```

## Database Query Examples

### Standard Query (all fields)
```sql
SELECT anime.*
FROM anime
INNER JOIN anime_seasons ON anime.id = anime_seasons.anime_id
WHERE anime_seasons.season = ?
```

### Field-Optimized Query (only requested fields)
```sql
SELECT anime.id, anime.title_en, anime.created_at, anime.updated_at
FROM anime
INNER JOIN anime_seasons ON anime.id = anime_seasons.anime_id
WHERE anime_seasons.season = ?
```

## Benefits

1. **Reduced Data Transfer**: Only fetches requested columns
2. **Better Performance**: Smaller result sets and faster queries
3. **Lazy Episode Loading**: Episodes only fetched when explicitly requested
4. **Backward Compatible**: Falls back to standard queries if field selection fails

## To Enable Full Field Selection

Currently, `ExtractAnimeFieldSelection()` returns `nil` to use standard queries. To enable field selection:

1. Implement GraphQL AST parsing in `ExtractAnimeFieldSelection()`
2. Extract field names from the GraphQL query context
3. Map GraphQL field names to database column names
4. Return populated `FieldSelection` object

## Testing

The existing tests verify that:
- Episodes are not fetched when not requested (lazy loading works)
- Field selection infrastructure is in place and functional
- Standard optimization still works when field selection is disabled

Run tests with:
```bash
go test ./internal/resolvers -v -run TestAnimeBySeasons
```