# WalletControl

WalletControl is a production-grade backend service designed for family budget management and automated financial analytics. 

## Architectural Design & Core Principles

The codebase enforces strict separation of concerns through concentric architectural layers. Dependencies point strictly inward, ensuring that the core business logic remains entirely decoupled from infrastructure details (databases, transport protocols, or loggers).

- **Domain (`internal/domain`)**: Contains pure business models (`Category`, `Transaction`) with zero external dependencies.
- **Usecase (`internal/usecase`)**: Implements application-specific business rules. It communicates with data sources exclusively through narrow interfaces (Interface Segregation Principle), enabling complete isolation and high-speed unit testing.
- **Repository (`internal/repository`)**: Handles persistence. It maps domain models to PostgreSQL utilizing `pgx/v5` for direct, high-performance database interaction and connection pooling.
- **Delivery (`internal/delivery/http`)**: Manages external communication. It leverages Go's built-in `net/http.ServeMux` for method- and path-based routing, accompanied by custom production-ready middleware.

---

## System Mechanics & Lifecycle

### 1. Request Lifecycle & Middleware Pipeline
When an HTTP request hits the server, it traverses a rigorously defined execution chain:
1. **Panic Recovery Middleware**: Traps any unhandled panics within downstream handlers. It captures the execution stack trace using `runtime/debug.Stack()`, logs the critical error via structured logging, and responds with a safe `500 Internal Server Error` to prevent process termination.
2. **Structured Logging Middleware**: Measures request duration using monotonic time (`time.Since`), injecting contextual attributes (HTTP method, request path, execution latency) into standard `log/slog` entries for downstream log aggregators.
3. **Mux Routing & Dispatch**: The standard `net/http.ServeMux` matches the request against registered path patterns and delegates execution to specific handlers.
4. **JSON Serialization & Usecase Execution**: Handlers decode incoming payloads, pass primitives or validated inputs down to the usecase layer, and serialize domain responses back into standard JSON streams.

### 2. Dependency Injection & Bootstrap Mechanics
The application eschews global state and reflection-based dependency injection containers. Initialization is handled explicitly in `cmd/app/main.go`:
- Configuration is parsed and validated from the environment via `godotenv`.
- A concurrent-safe connection pool (`pgxpool.Pool`) is initialized, with immediate health verification executed via `pool.Ping(ctx)`.
- Repositories are instantiated with the database pool, injected into the usecases, and bound to the HTTP router, forming a clean, traceable dependency graph.

### 3. Asynchronous Background Processing
Financial analytics aggregation runs asynchronously via a dedicated background worker (`internal/worker/analytics.go`):
- The worker executes inside an independent goroutine, triggered periodically by a `time.Ticker`.
- It performs rolling 30-day aggregations (`GROUP BY` queries on timestamps) independently of the request-response cycle.
- **Context-Aware Lifecycle**: The worker loop listens to `ctx.Done()`. Upon application shutdown, the context is cancelled, ensuring that the worker terminates cleanly without leaking goroutines or executing orphaned queries.

### 4. Graceful Shutdown Protocol
To prevent data corruption and dropped client requests during deployments or restarts, the server implements coordinated shutdown:
- The HTTP server runs in a separate background goroutine via `server.ListenAndServe()`.
- The main OS thread blocks on a buffered channel listening for `SIGINT` and `SIGTERM` signals.
- Upon signal reception, a timeout-bounded context (`context.WithTimeout` set to 5 seconds) is dispatched to `server.Shutdown(ctx)`. This drains active HTTP connections, finishes ongoing request cycles, closes the database pool safely, and exits.

---

## API Endpoints

- `POST /categories` — Creates a new expense category.
- `GET /categories` — Retrieves all configured categories.
- `POST /expense` — Records a new financial transaction linked to a category.
- `GET /analytics` — Fetches aggregated spending statistics for the last 30 days.

---

## Getting Started

### Prerequisites
- Docker & Docker Compose
- Make