# 📝 Tasks API — A Go Standard Library REST API

[![Go Version](https://img.shields.io/badge/Go-1.22%2B-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org)
[![Database](https://img.shields.io/badge/PostgreSQL-15%2B-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)](https://www.postgresql.org/)

A RESTful API for task management, built entirely with **Go's standard library (`net/http`)** — zero third-party web frameworks.

Built as part of a structured backend transition from **PHP to Go**, prioritizing idiomatic Go patterns and explicit design decisions over framework magic.

---

## 🎯 Motivation & Engineering Goals

Moving from PHP frameworks (Laravel) to Go requires a mindset shift from implicit convention to explicit design. This project was built to practice key Go backend patterns beyond tutorial-level CRUD:

- **Explicit dependency injection**: handlers are methods on a `*TaskStore` struct, not standalone functions relying on global state.
- **Idiomatic HTTP routing**: Go 1.22+ method-aware `http.ServeMux` patterns (`GET /tasks/{id}`), no router library.
- **SQL injection prevention**: every query uses parameterized statements (`$1`, `$2`...) via `database/sql`.
- **Deliberate HTTP status mapping**: `400` for invalid input, `404` for missing resources, `500` for unexpected failures — handled as distinct cases, not lumped together.
- **Table-driven tests**: handler behavior (status codes, response bodies) and validation logic tested in isolation with `net/http/httptest`, no real server needed.

---

## 🛠 Tech Stack

| Component | Tech / Library | Reason |
| :--- | :--- | :--- |
| **Language** | Go 1.22+ | Native performance, static typing, built-in concurrency primitives |
| **HTTP Routing** | Standard library `net/http` | Method-aware routing and path parameters, no framework needed |
| **Database** | PostgreSQL 15+ | Relational storage for task persistence |
| **DB Driver** | `jackc/pgx/v5` | High-performance PostgreSQL driver for `database/sql` |
| **Config** | `joho/godotenv` | Environment variable loading from `.env` for local development |

---

## 📁 Project Architecture

Following the standard Go `cmd/` + `internal/` project layout:

```text
tasks-api/
├── cmd/
│   └── tasks-api/
│       └── main.go           # Application entrypoint
├── internal/
│   └── tasks/
│       ├── handlers.go       # HTTP handlers as methods on *TaskStore
│       ├── handlers_test.go  # Table-driven tests for handlers and validation
│       ├── models.go         # Task struct and validation logic
│       ├── routes.go         # Route registration
│       └── store.go          # Database access layer (SQL queries)
├── .env.example               # Sample environment configuration
├── schema.sql                  # Table definition for PostgreSQL
├── go.mod
└── README.md
```

---

## 🔑 Key Decisions

- **Dependency injection over globals**: earlier versions of this project used package-level variables for storage. That made state implicit and forced tests to manually reset it before every run. Refactoring to inject `*TaskStore` explicitly fixed both problems.
- **Parameterized queries everywhere**: no string concatenation into SQL, ever — every value is passed positionally to prevent injection.
- **Standard library over a framework**: routing, JSON encoding/decoding, and HTTP handling are all built directly on `net/http` and `encoding/json`, to understand what a framework like Gin or Echo would otherwise abstract away.

---

## 🚀 Running Locally

1. Clone the repo
2. Copy `.env.example` to `.env` and fill in your database URL
3. Create the database and run `schema.sql` against it
4. `go run ./cmd/tasks-api`

## ✅ Tests

```bash
go test ./...
```

Covers handler status codes and response bodies (including error paths: malformed JSON, invalid IDs, failed validation) and the validation logic in isolation, using table-driven tests.

---

## 🔭 What I'd Do Differently With More Time

- Add synchronization (or move fully to the database as the single source of truth with proper transaction boundaries) — the current version isn't yet protected against race conditions under concurrent writes. This is addressed in the next phase of my learning path.
- Graceful shutdown on `SIGTERM`/`SIGINT`, so in-flight requests aren't dropped during a deploy.
- Structured logging and basic observability (currently the app has neither).
