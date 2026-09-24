# test-project-2

A todo app with a Go backend (standard library + [pgx](https://github.com/jackc/pgx)) and PostgreSQL 17. The HTML/CSS/JS frontend is embedded in the binary, so the whole app ships as one container.

- **Preview:** https://abandon-above-toast.preview.basable.com (deployed from `dev`)
- **Live:** https://abandon-above-toast.live.basable.com (deployed from `main`)

## Features

- Add, complete, edit (double-click) and delete tasks
- Filter by All / Active / Done, with a progress ring showing how many are done
- **Automatic cleanup:** when there are more than 100 todos, the oldest 50 are deleted

## API

| Method   | Path              | Body                               | Response          |
|----------|-------------------|------------------------------------|-------------------|
| `GET`    | `/api/todos`      | –                                  | `200` todo list   |
| `POST`   | `/api/todos`      | `{"title": "..."}`                 | `201` new todo    |
| `PATCH`  | `/api/todos/{id}` | `{"title": "...", "done": true}`   | `200` updated todo |
| `DELETE` | `/api/todos/{id}` | –                                  | `204`             |

Titles must be 1–500 characters. Both `PATCH` fields are optional.

| Path      | Purpose                              |
|-----------|--------------------------------------|
| `/health` | Liveness: always `200 OK`            |
| `/ready`  | Readiness: `200` once the database responds |

## Configuration

| Variable       | Required | Default | Description                  |
|----------------|----------|---------|------------------------------|
| `DATABASE_URL` | yes      | –       | PostgreSQL connection string |
| `PORT`         | no       | `8080`  | HTTP listen port             |

The `todos` table is created automatically when the app starts.

## Running locally

```sh
docker run -d --name todos-db -e POSTGRES_PASSWORD=pw -p 5432:5432 postgres:17-alpine
DATABASE_URL='postgres://postgres:pw@localhost:5432/postgres?sslmode=disable' go run .
```

Then open http://localhost:8080.

## Tests

```sh
go test ./...
```

The tests use an in-memory store, so they don't need a database.

## Project layout

| Path              | Contents                                            |
|-------------------|-----------------------------------------------------|
| `main.go`         | Startup, config, graceful shutdown                  |
| `server.go`       | Routes and embedded static files                    |
| `handlers.go`     | HTTP handlers, validation, cleanup limits           |
| `store.go`        | `Todo` model and `Store` interface                  |
| `pgstore.go`      | PostgreSQL store                                    |
| `pgtrim.go`       | Oldest-first cleanup (uses an advisory lock so it runs one at a time) |
| `static/`         | Frontend (HTML, CSS, JS)                            |
| `k8s/`            | Kubernetes manifests                                |
| `.basable/`       | Platform sizing (`config.yaml`) and build script (`verify.sh`) |

## Deployment

When you push to `dev`, GitHub Actions builds the image and the app deploys to preview. When you merge `dev` into `main`, it deploys to live.

`k8s/` contains:

- the app `Deployment`, `Service` and `HTTPRoute`
- a single-replica PostgreSQL `StatefulSet` with a 5Gi volume
- the database credentials, in `k8s/postgres-secret.yaml`
