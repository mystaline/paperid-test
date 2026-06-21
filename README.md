# paperid-test

disbursement service, single endpoint `POST /disburse`.

requires: docker + go 1.26

## setup

1. copy env, then start the service (server & postgres):

   ```sh
   cp .env.example .env
   docker compose up -d --build
   ```

   note: server runs on `:8080`, postgres on `:55432` to avoid collision with existing postgres instance in reviewer/user machine, schema and seed load automatically on first run.

   docker is the simple path, it serves postgres and loads schema & seed for reviewer/user, so no manual db setup. the steps below are only if you docker not available.

   without docker (assumes postgres is already running):

   ```sh
   createdb paperid_test
   psql -d paperid_test -f pkg/db/init/01_schema.sql -f pkg/db/init/02_seed.sql
   go run .   # set POSTGRES_HOST/PORT/USER/PASSWORD/DB env if your db differs from localhost:5432 postgres/postgres
   ```

   or directly create db "paperid_test" in any db client then manually execute queries in pkg/db/init (01_schema.sql and 02_seed.sql)

2. hit the endpoint — user `100000001` is seeded with balance:

   ```http
   POST http://localhost:8080/disburse
   Content-Type: application/json

   {"amount":50000,"userId":"100000001"}
   ```

   or via curl:

   ```sh
   curl http://localhost:8080/disburse -H 'content-type: application/json' -d '{"amount":50000,"userId":"100000001"}'
   ```

   note: paste the raw request above into postman/bruno/insomnia, or run the curl.

   note: since im using BIGINT (snowflake id), input/output related to ids is string to prevent loss precision.

## tests

unit test without db, repos are mocked, just to ensure core logic is verified:

```sh
go test ./... -count=1
```

integration test needs postgres up, runs against the real db (atomic debit, rollback on log failure, concurrent overdraft), to ensure i made battle-tested service:

```sh
docker compose up -d
go test ./... -count=1
```
