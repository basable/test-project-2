# test-project-2

A todo app: Go (standard library + pgx) backend with an embedded HTML/JS frontend, backed by PostgreSQL 17.

## API

| Method | Path              | Body                            |
|--------|-------------------|---------------------------------|
| GET    | `/api/todos`      | –                               |
| POST   | `/api/todos`      | `{"title": "..."}`              |
| PATCH  | `/api/todos/{id}` | `{"title": "...", "done": true}` |
| DELETE | `/api/todos/{id}` | –                               |

`/health` is liveness, `/ready` checks the database.

## Local run

```sh
docker run -d -e POSTGRES_PASSWORD=pw -p 5432:5432 postgres:17-alpine
DATABASE_URL=postgres://postgres:pw@localhost:5432/postgres?sslmode=disable go run .
```

## Deployment

`k8s/` has the app Deployment and a single-replica Postgres StatefulSet (5Gi volume). Credentials are in `k8s/postgres-secret.yaml`.
