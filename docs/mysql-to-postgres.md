# Moving the read store from MySQL to Postgres

Plan for retiring PlanetScale/Vitess in favour of Postgres on AWS RDS.
Written 2026-08-22, before implementation.

## Why

Stated preference is Postgres, and migrations are materially easier there.
Concrete examples from this codebase:

- **No foreign keys.** `000032_create_tags_table` says so in its own comment:
  "Foreign keys not used due to Vitess/PlanetScale compatibility."
- **`migrate up` cannot build a database from scratch.** An earlier migration
  uses `DELIMITER`, which golang-migrate cannot parse, so the chain dies
  partway. New environments cannot be bootstrapped from migrations alone.
- **Accent folding by hand.** `000038_add_url_slug_to_anime_staff` needed a
  ~40-branch nested `REPLACE` chain to fold accents inside a generated column.
  Postgres does that with `unaccent()`.

## The prize is bigger than migrations

Postgres is **already the source of truth**. The entire CDC pipeline exists to
translate Postgres into MySQL:

```
scraper Postgres -> Debezium -> Kafka -> sync services -> MySQL -> anime-api
```

Same engine on both ends removes the translation, and with it:

```
8   sync deployments in prod
16  kafka topics
13  confluent pods
1   weekly reconcile job (db-sync-postgres-prod-to-staging)
    Debezium connectors and their replication slots
```

That is not incidental. This week that pipeline failed silently for a month,
left 18,311 rows missing from staging while every component reported healthy
and consumer lag sat at zero, and it still cannot propagate deletes — 2,405
anime exist in staging MySQL that prod deleted. Same-engine replication is
native, and none of that machinery has to be maintained or monitored.

## Most of the data does not need migrating

This is the point that makes the cutover small. Split the tables by who owns
the truth:

**Derived from Postgres — rebuild, do not migrate.** Everything here arrived
via CDC and can be regenerated from the source of truth:

```
anime_character_staff_link  427k rows   276 MB
episodes                    149k        196 MB
anime                        25k        110 MB
anime_character             216k         85 MB
anime_staff                  21k         13 MB
anime_tags                   33k         12 MB
anime_seasons, anime_news, anime_fanart, anime_streaming_platform, tags
```

**Authoritative in MySQL — must be migrated.** Nothing upstream can rebuild
these:

```
sessions          28,586     ephemeral, can be dropped (users re-login)
refresh_tokens    28,421     ephemeral, can be dropped
credentials        2,387
user_anime         2,309
users                102
user_list              0
password_resets        0
```

Discounting the two ephemeral tables, **the irreplaceable dataset is roughly
4,800 rows.** That is the entire risk surface of the data migration. The 710 MB
catalogue is derived and disposable.

Treat those as two separate migrations with different risk profiles. The
catalogue is a rebuild that can be re-run as often as needed; the user data is
a one-way move that must be exact.

## Starting position on AWS

Verified, because "already set up" is only half true:

| repo | state |
|---|---|
| `aws-root-tf` | organisations, identity center, GitHub OIDC, backend. Account scaffolding. |
| `vpc-tf` | VPC definitions present |
| `rds-tf` | **reusable module** for Aurora PostgreSQL and standard RDS, Postgres 16.4. Examples only. |
| `terraform-storage` | **empty** — one README, initial commit |
| — | **no database provisioned** |

`terraform-storage` is where `rds-tf` should be instantiated. The module is
written and takes `vpc_id`, `subnet_ids`, `allowed_cidr_blocks`, instance class
and `instance_count` (writer + readers), so the remaining work there is an
instantiation and its inputs, not new module code.

## Target shape

```
RDS Postgres primary        scraper writes
        |
        | native read replica (RDS-managed)
        v
RDS Postgres replica        anime-api reads
```

No Debezium, no Kafka, no sync services. RDS read replicas are managed, so the
read/write isolation the current design buys with a whole pipeline comes for
free.

**Open question: where does the scraper's Postgres live?** Either it moves to
RDS as the primary, or it stays in-cluster and logically replicates to RDS.
Moving it is simpler and removes another in-cluster stateful workload; keeping
it avoids putting a batch scraper's write load on the same instance the site
reads from. Worth deciding before provisioning, since it changes the instance
sizing.

## What changes in anime-api

- **Driver and dialect** — GORM `mysql` to `postgres`.
- **Migrations.** 38 MySQL migrations do not translate one for one, and one of
  them cannot run from scratch anyway. Squash to a single baseline schema for
  the new database rather than porting the history; the history's value is the
  reasoning in its comments, which the baseline can carry forward.
- **Dialect-specific SQL to rewrite:**
  - `000038` generated column: MySQL `REGEXP_REPLACE` + `REPLACE` chain becomes
    `unaccent()`
  - `000032` tags migration: `JSON_TABLE` becomes `jsonb_array_elements_text`
  - `ON DUPLICATE KEY UPDATE` becomes `ON CONFLICT ... DO UPDATE`
  - `anime_episode_count` triggers become plpgsql
  - Backtick quoting becomes double quotes
- **Foreign keys become available again** — the constraints skipped for Vitess
  compatibility can be real.
- **`interpolateParams`, `multiStatements`** and the other MySQL DSN flags go
  away.

## Phases

Ordered so nothing is bet before it is proven.

1. **Provision.** Instantiate `rds-tf` from `terraform-storage`. Decide
   primary-only vs primary + replica, and sizing. Nothing else depends on this
   being right first time — RDS instance classes are changeable.

2. **Baseline schema.** Author the Postgres schema for anime-api as one
   migration, carrying forward the reasoning from the 38 it replaces.

3. **Rebuild the catalogue into it** from the scraper's Postgres. This is a
   rehearsal, not a cutover — it can be run repeatedly and thrown away, and it
   is how the schema gets validated against real data.

4. **Point staging anime-api at RDS.** Run both read stores in parallel and
   compare: row counts per table, and the same GraphQL queries against both.

5. **Migrate the user data** — `users`, `credentials`, `user_anime`,
   `user_list`. Roughly 4,800 rows, and the only part that is irreversible.
   `sessions` and `refresh_tokens` are deliberately not migrated; the cost is
   that everyone signs in again.

6. **Cut prod over.** DBHOST change per service.

7. **Retire the pipeline.** Debezium connectors, Kafka topics, the 8 sync
   deployments, the reconcile CronJob. Last, and only once prod has been stable
   on Postgres long enough to be sure.

## Risks worth naming

- **PlanetScale is managed and has not caused trouble.** It is one of the more
  reliable things in this stack. RDS is also managed, so this swaps one managed
  service for another — the risk is in the move, not the destination. That is
  materially different from self-hosting Postgres in-cluster, which is not what
  this plan proposes.
- **Egress and latency.** anime-api runs in-cluster and would read across the
  network to AWS. Needs measuring against current PlanetScale latency, which is
  also remote, so this may be a wash — but measure rather than assume.
- **Cost.** PlanetScale's bill is replaced by RDS instance plus storage plus
  cross-AZ and egress traffic. The pipeline retirement claws some back (13
  confluent pods, 8 sync deployments), but it is a different cost shape.
- **Step 7 is the point of no return.** Once Debezium and the sync services are
  gone, falling back to MySQL means rebuilding them. Keep them running,
  idle-but-working, through at least one full scrape cycle after cutover.

## Related

- `anime_scraper/docs/manga-and-works.md` — the manga/works plan, which is
  sequenced around this one. New entity families built before the move pay the
  CDC pipeline cost twice.
