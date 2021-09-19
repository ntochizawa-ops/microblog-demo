microblog
=========

## Environment variables

| Variable | Type | Required | Description |
|---|---|---|---|
| `DATABASE` | string | yes | Database name (`projects/PROJECT_ID/instances/INSTANCE_ID/databases/DATABASE_ID`) |
| `PORT` | string | no | default: `"8080"` |
| `LOG_LEVEL` | string | no | default: `"info"` |

## Develop with Spanner emulator

Install Cloud Spanner emulator.

```bash
gcloud components install cloud-spanner-emulator
```

Start Cloud Spanner emulator.

```bash
gcloud emulators spanner start
$(gcloud emulators spanner env-init)
```

If you haven't your instance yet, create a new instance and create a new database with [hammer](https://github.com/daichirata/hammer).

```bash
gcloud config configurations create emulator
gcloud config set auth/disable_credentials true
gcloud config set api_endpoint_overrides/spanner http://localhost:9020/
gcloud spanner instances create microblog \
  --project microblog \
  --config emulator-config \
  --description "Microblog" \
  --nodes 1
gcloud config configurations activate default
hammer create spanner://projects/microblog/instances/microblog/databases/microblog schema.sql
```

Run app with the database on the emulator.

```bash
DATABASE=projects/microblog/instances/microblog/databases/microblog go run *.go
```
