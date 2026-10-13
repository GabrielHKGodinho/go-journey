# Case Study: Uber and WhatsApp vs. My Scaling Design

Week 10 — System Design. Compares my scaling design for `investment-engine`
(peak ~25 writes/s, ~600 reads/s, ~30 GB/year) with two real architectures.

## Comparison

| Dimension | My design | Uber | WhatsApp | What they did differently, and why |
|---|---|---|---|---|
| 1. Where state lives | Source of truth in Postgres; Redis is only a read cache (a derived copy) | Geofence service: a full copy of every polygon in each instance's process memory; the source of truth stays in the datastores it polls | Messages in transit live in Mnesia (in-memory, ~2 TB RAM across 16 partitions); history lives on the user's phone | I guessed an external cache. Uber keeps the data inside the process because each lookup is CPU-heavy and runs ~170k times per second, so a network hop per query would waste latency. It works only because the dataset fits in RAM and changes rarely. |
| 2. How it scales | Stateless tier horizontal (api ×3, competing consumers, priceservice); stateful tier vertical with replicas for availability (Postgres primary + sync standby, RabbitMQ ×3 quorum) | Geofence: ~40 identical stateless instances (horizontal). Storage: sharding (Schemaless on MySQL) once one database was not enough | Vertical first: tuned Erlang/FreeBSD to hold ~2M connections per server; added machines when needed (>11,000 cores by 2016) | I did not expect a giant to lean on vertical scaling. Rule: scale stateless layers out (cheap), scale stateful layers up for as long as possible (sharding is expensive). |
| 3. Consistency between copies | Sync standby (WAL), quorum queues; idempotent consumer handles duplicates | Eventual consistency: every instance refreshes on the same deterministic schedule, so copies differ only for the rebuild time plus clock skew | Not well documented publicly; the 2014 write-up itself warns little is known about the overall architecture | I expected periodic polling to be too expensive. It is cheap compared to the query load; the hidden cost is that all instances hit the sources at the same moment instead of spreading calls with jitter. |
| 4. Failure handling / retry | Separate retry queue with a fixed 5 s TTL, `x-retry-count`, max 3 attempts, permanent errors to the DLQ | Ladder of retry topics with increasing delays, DLQ at the end, tooling to list/purge/replay the DLQ, code bugs straight to the DLQ | "Let it crash": a failing Erlang process dies alone and a supervisor restarts it with clean state | Same core idea (publish a copy, then ack: success = a conclusive outcome). Differences I lack: increasing delays (a 60 s outage sends ~83% of orders to my DLQ), a DLQ replay tool, and protection against panics, which bypass my retry logic entirely and create a poison-message crash loop. |
| 5. Data kept forever? | Yes: orders and executions are financial records; the portfolio is derived from them | Geofence copies are disposable; the source of truth lives in the datastores | No: messages leave the server after delivery, media is not kept long-term | WhatsApp removed the hardest problem (storing everything forever) with a product decision. State grows with messages in transit, not with history. My state grows with history, so it cannot do the same. |
| 6. Language / runtime, and why | Go: goroutines (not OS threads), static typing, single binary. Honestly also chosen for my career goal | Go was chosen for this one service: CPU-heavy lookups plus background refresh would block Node.js's single thread; the rest of Uber was still mostly Node.js and Python | Erlang on the BEAM: one lightweight process per connection, per-process isolation, supervisors, hot code swapping | Uber picked the language per service, by workload, and limited the risk to one small new service. An idle goroutine costs ~2.5 KB (measured), the same order as an Erlang process; the real difference is failure isolation. |
| 7. Team size × number of services | 1 developer × 3 Go services (7 counting Postgres, RabbitMQ, Jaeger) | ~2,200 critical microservices; needed a whole architecture (DOMA: domains, layers, gateways, extensions) to manage them | ~10 people on the Erlang system, doing both development and operations | The number of services follows the number of teams (Conway's law). Uber itself says startups can hold off on microservices. My 3 services are a demonstration; a one-person startup should start with a modular monolith. |

## Takeaways

1. Copy the reasoning, not the decision. Each choice above is tied to a scale, a team and a date.
2. The same idea can be cheap in one tool and expensive in another (retry count in a header: costly in Kafka, cheap in RabbitMQ thanks to queue TTL).
3. My safety nets cover returned errors, not panics. Fix: recover in the handler and treat a panic as a permanent failure; quorum queues add a broker-side delivery limit.
4. Interview-rubric answers ("great" vs "good") grade the interview's problem, not mine. A conditional `UPDATE ... WHERE status = 'PENDING'` beats a Redis TTL lock for a money system at my scale.

## Sources

- Uber — How We Built Uber Engineering's Highest Query per Second Service Using Go (2016): https://www.uber.com/blog/go-geofence-highest-query-per-second-service
- Uber — Building Reliable Reprocessing and Dead Letter Queues with Apache Kafka (2018): https://www.uber.com/blog/reliable-reprocessing/
- Uber — Introducing Domain-Oriented Microservice Architecture: https://www.uber.com/blog/microservice-architecture/
- Uber — Why Uber Engineering Switched from Postgres to MySQL (2016): https://www.uber.com/blog/postgres-to-mysql-migration/
- InfoQ — Uber Engineering Moving from Postgres to MySQL: https://www.infoq.com/news/2016/08/Uber-Engineering-Postgres-MySQL
- Hello Interview — Design Uber: https://www.hellointerview.com/learn/system-design/problem-breakdowns/uber
- HighScalability — The WhatsApp Architecture Facebook Bought For $19 Billion (2014): https://highscalability.com/the-whatsapp-architecture-facebook-bought-for-19-billion/
- HighScalability — How WhatsApp Grew to Nearly 500 Million Users, 11,000 cores, and 70 Million Messages a Second (2016): https://highscalability.com/how-whatsapp-grew-to-nearly-500-million-users-11000-cores-an/
- Martin Kleppmann — How to do distributed locking (2016): https://martin.kleppmann.com/2016/02/08/how-to-do-distributed-locking.html
