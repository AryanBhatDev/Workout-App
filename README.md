# Workout App

HTTP API for registering users, issuing auth tokens, and storing workouts with exercise entries. Postgres is the database. Migrations run automatically when the server starts.

## Requirements

- Go 1.26 or newer
- Docker, for the Postgres containers in `docker-compose.yml`

## Run

Start Postgres:

```bash
docker compose up -d db
```

The app database listens on `localhost:5432` with user `postgres`, password `postgres`, and database `postgres`.

Start the API (default port `8080`):

```bash
go run .
```

Override the port or the connection string:

```bash
go run . -port 8080 -"db conn string" "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"
```

Check that the process is up:

```bash
curl http://localhost:8080/health
```

## Auth

Workout routes require `Authorization: Bearer <token>`.

Register:

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"username":"ada","email":"ada@example.com","password":"secret","bio":"lifting"}'
```

Create a token (valid for 24 hours):

```bash
curl -X POST http://localhost:8080/tokens/authentication \
  -H "Content-Type: application/json" \
  -d '{"username":"ada","password":"secret"}'
```

The response field is `auth_token`. Send that value on later requests:

```bash
-H "Authorization: Bearer <auth_token>"
```

## Routes

| Method | Path | Auth | What it does |
|---|---|---|---|
| `GET` | `/health` | no | Liveness check |
| `POST` | `/users` | no | Register a user. Returns `201` and `{"user": ...}` |
| `POST` | `/tokens/authentication` | no | Check username and password. Returns `201` and `{"auth_token": ...}` |
| `POST` | `/workouts` | yes | Create a workout and its entries. Returns `200` and `{"workout": ...}` |
| `GET` | `/workouts/{id}` | yes | Load one workout and its entries |
| `PUT` | `/workouts/{id}` | yes | Registered on the get handler |
| `DELETE` | `/workouts/{id}` | yes | Delete a workout. Returns `204` |

Create a workout:

```bash
curl -X POST http://localhost:8080/workouts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <auth_token>" \
  -d '{
    "title": "Morning Strength Training",
    "description": "Full body strength training session",
    "duration_minutes": 60,
    "calories_burned": 500,
    "entries": [
      {
        "exercise_name": "Squats",
        "sets": 3,
        "reps": 12,
        "weight": 100.5,
        "notes": "Felt strong today",
        "order_index": 1
      }
    ]
  }'
```

Errors are JSON objects with an `error` string. A bad id is `400`, a missing workout is `404`, and a database failure is `500`.

## Tests

Start the test database, which is published on port `5433`:

```bash
docker compose up -d test_db
go test ./internal/store/
```

`setupTestDB` migrates that database and truncates `users`, `workouts`, and `workout_entries` before each test.

## Layout

```text
main.go                         start the HTTP server
internal/app/app.go             open the DB, run migrations, build handlers
internal/routes/routes.go       URL to handler
internal/middleware/            bearer token auth
internal/api/                   HTTP handlers
internal/store/                 SQL
internal/tokens/                token generation
internal/utils/                 JSON responses and id parsing
migrations/                     goose SQL, applied on startup
```

Request path:

```text
routes.Setup
  -> handler method
    -> store method
      -> Postgres
```
