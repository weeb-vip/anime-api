# Migrations

## Numbering starts again at 000040

Versions **37, 38 and 39 are permanently retired**. They created `anime_news` and
`anime_fanart` and added the news language/reference columns; that schema now belongs to
**news-ingest**, which owns those tables and migrates them under its own migration table
(`__migrations_news-ingest`).

**Do not reuse 37–39 for anything else.** Every database that ran them still records having
done so, and golang-migrate identifies a migration solely by its number. A new file numbered
37 would be silently skipped on those databases and applied on fresh ones — the two would
diverge with nothing to indicate it. Leaving the numbers unused costs nothing; reusing them
costs a schema drift that only shows up in production.

The next migration added here should be `000040_*`.

## Existing databases record version 39

Since 37–39 no longer exist in this repository, a database still at version 39 has a version
this source cannot explain, and `migrate up` will fail looking for a file that is gone.

Each environment needs its recorded version set back to 36 **once**, at the same time this
change is deployed:

```sql
UPDATE schema_migrations SET version = 36, dirty = 0;
```

That is safe because 37–39 only ever created news tables, and those tables are deliberately
left in place — news-ingest adopts them. Nothing is dropped and no data moves; only the
bookkeeping changes, so that anime-api's history matches the files it actually ships.

Verify afterwards with `SELECT * FROM schema_migrations;` — it should read `36 | 0`, and
`migrate up` should then be a no-op until 000040 exists.
