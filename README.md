# Practicum Score API

[![CI](https://github.com/Dwiyuda/practicum-score-graphql/actions/workflows/ci.yml/badge.svg)](https://github.com/Dwiyuda/practicum-score-graphql/actions/workflows/ci.yml)

**[Try the live demo](https://dwiyuda.github.io/practicum-score-graphql/)**: the real schema and resolvers,
compiled to WebAssembly and running in your browser. It has no server or PostgreSQL behind it: data lives
in the in-memory store and resets on refresh. Use Docker Compose below for the full stack.

A GraphQL API in Go for recording practicum (lab) scores and attendance, backed by PostgreSQL.
It answers the questions a teaching assistant asks every week: how is each student doing, who is
missing sessions, and how does a whole class compare.

Built with [gqlgen](https://gqlgen.com) (schema first), [pgx](https://github.com/jackc/pgx), and
Docker Compose. All demo data is synthetic; names are fictional.

## Run it

Requires Docker.

```bash
docker compose up --build
```

Open <http://localhost:8080> for the GraphQL playground. The API container runs the migrations and
seeds demo data on start.

Example query:

```graphql
{
  classSummary(className: "A") { students average attendanceRate }
  students(className: "B") { nim name average attendanceRate }
}
```

Example mutation (creates the score or replaces the existing one for that student and session):

```graphql
mutation {
  recordScore(input: { studentId: 1, sessionId: 2, value: 87.5, attended: true }) {
    value
    student { name average }
  }
}
```

## Schema

| Type | Fields |
|---|---|
| `Student` | `nim`, `name`, `className`, `scores`, `average`, `attendanceRate` |
| `Session` | `number`, `topic`, `heldOn` (YYYY-MM-DD) |
| `Score` | `student`, `session`, `value` (0–100), `attended` |
| `ClassSummary` | `students`, `average`, `attendanceRate` |

Queries: `students(className)`, `student(id)`, `sessions`, `classSummary(className)`.
Mutations: `createStudent`, `createSession`, `recordScore`.

## Layout

```
cmd/server/            HTTP server, startup migrations
cmd/wasm/              WebAssembly entry point for the browser demo
demo/                  Static page for the browser demo, published to GitHub Pages
graph/                 GraphQL schema, gqlgen generated code, resolvers
internal/grading/      score rules: validation, average, attendance rate
internal/store/        Store interface, PostgreSQL implementation, SQL migrations, in-memory store, demo seed
```

Score rules live in `internal/grading`, away from storage, so they are tested without a database.
The database enforces the same limits (`CHECK (value BETWEEN 0 AND 100)`, one score per student and
session) so bad data cannot get in through another path.

## Tests

```bash
go test ./...                       # unit + resolver tests (in-memory store)
DATABASE_URL=postgres://... go test ./internal/store/   # integration test against a real PostgreSQL
```

GitHub Actions runs `go vet`, the full test suite with the race detector against a PostgreSQL
service container, checks that the generated GraphQL code is up to date, and builds the Docker image.

## Notes

- Migrations are plain SQL files embedded into the binary and applied in order at startup; each is
  recorded in `schema_migrations`, so restarting is safe.
- `-seed` only inserts data when the database has no students yet.
- Hosting the server with PostgreSQL is optional; see [DEPLOY.md](DEPLOY.md).
- I used an AI coding agent (Claude Code) while building this, and checked its output with the tests above.

Built by Dwi Yuda.
