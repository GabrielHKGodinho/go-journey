# Go Journey

Daily practice log for my transition from PHP to Go — following a
structured 5-month roadmap (July–December 2026), covering language
fundamentals, concurrency, and distributed systems.

This is not a single project — it's a collection of small, focused
exercises, each in its own folder, tracking what I learn day by day.
For a complete, standalone project built along the way, see
[tasks-api](https://github.com/GabrielHKGodinho/tasks-api) — a REST
API with PostgreSQL, dependency injection, and tests.

## Structure

Organized by week and topic, roughly following the order I learned
things in:

```text
go-journey/
├── week0-basics/          # syntax, defer, slices, closures
├── week1-structs/         # structs, methods, receivers, interfaces
├── week2-http/            # net/http, JSON, errors (tasks-api lives
│                            here locally, but is tracked separately)
└── week3-concurrency/     # goroutines, channels, sync, context
```

Each exercise folder is an independent Go module (`go.mod`), since
most of them aren't meant to import from each other — they're
standalone practice, not a single application.

## Why this exists

Beyond the learning itself, this repo is meant to show consistency —
one exercise at a time, most days, rather than a single burst of
effort. It's the "show your work" counterpart to the more polished
`tasks-api` project.

## Notes

Some exercises are minimal on purpose (a single concept, isolated).
Others — like the worker pool and the RWMutex cache under
`week3-concurrency/` — are closer to small standalone projects,
built to practice a specific pattern end to end.