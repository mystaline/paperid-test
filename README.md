# paperid-test

disbursement service, single endpoint `POST /disburse`.

## setup

1. start the service (server & postgres):

   ```sh
   docker compose up -d --build
   ```

   note: server runs on `:8080`, postgres on `:55432` to avoid collision with existing postgres instance in reviewer machine, schema and seed load automatically on first run.

2. hit the endpoint — user `100000001` is seeded with balance:

   ```sh
   curl http://localhost:8080/disburse -H 'content-type: application/json' -d '{"amount":50000,"userId":"100000001"}'
   ```

   note: or use postman/any http client, hit request with method POST to `http://localhost:8080/disburse` with that json body.

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
