# Go Backend Template – AI Agent Rules

This project is a production-grade Go backend template built as a modular monolith using Chi, Ent, Redis, Postgres, and observability tools (Sentry, PostHog).

AI agents contributing to this codebase MUST follow these rules strictly to maintain consistency, scalability, and production readiness.

---

# 1. Architecture Principles

- This is a **modular monolith**, not microservices.
- Each feature must live inside `internal/modules/<module-name>`.
- Modules must be **independent and isolated** (no cross-module imports unless via interfaces).
- Business logic MUST NOT live in handlers or middleware.
- Handlers are thin: request parsing → service call → response.

---

# 2. Dependency Rules

- `internal/app` is the entry point (composition root).
- `server` depends on `container`, NOT the other way around.
- All external services (DB, cache, analytics, etc.) must be injected via `container`.
- NEVER use global variables for stateful services.
- Interfaces are preferred for:
  - cache
  - storage
  - mailer
  - external APIs

---

# 3. Database (Ent) Rules

- All DB access must use Ent (`ent.Client`).
- No raw SQL unless explicitly justified.
- Queries must be performed in the service layer, not handlers.
- Always handle `NotFound` errors explicitly.
- Avoid returning Ent entities directly to API responses.

Use DTOs when:
- caching data
- exposing API responses
- data contains sensitive fields

---

# 4. Caching (Redis) Rules

- Redis must always be accessed via `cache.Cache` interface.
- Never call Redis client directly in business logic.
- Use cache-aside pattern:
  1. check cache
  2. fallback to DB
  3. write back to cache

- Cache keys must be namespaced:
  - `user:{id}`
  - `post:{id}`
- Cached data should be JSON-encoded DTOs, not raw Ent structs.

---

# 5. Authentication & Authorization

- JWT is used for authentication.
- Claims must include:
  - user_id
  - role
  - permissions

- Middleware must:
  - validate token
  - attach claims to request context

- Authorization must use:
  `RequirePermission(mode, permissions...)`

Rules:
- `super_admin` bypasses all permission checks
- Permissions must be string constants in `internal/permissions`

---

# 6. Middleware Rules

- Middleware must be stateless.
- Order of middleware matters:
  1. Request ID
  2. Recoverer
  3. Logger
  4. Security
  5. CORS
  6. Rate limiting
  7. Auth
  8. Authorization

- Middleware must NOT contain business logic.
- Middleware must NOT call DB directly.

---

# 7. Error Handling Rules

- Use centralized error types (`internal/errors`).
- Do NOT use `http.Error` in business logic.
- Errors must be:
  - typed
  - consistent
  - reusable

- Handlers should return structured responses, not raw errors.

---

# 8. Response Rules

- All API responses must use `internal/response`.
- Responses must be consistent:
  - success
  - error
  - created
  - no content

- Never return raw structs directly without response wrapper.

---

# 9. Logging & Observability

- All logs must go through structured logger (zap).
- Logs must include:
  - request ID
  - method
  - path
  - latency
  - status code

- Sentry is used for error tracking:
  - Only log unexpected/system errors to Sentry
  - Do NOT spam Sentry with expected validation errors

- PostHog is used for analytics:
  - Only track meaningful business events
  - Never track sensitive data (passwords, tokens, PII unless anonymized)

---

# 10. Validation Rules

- All incoming request payloads must be validated.
- Use `go-playground/validator`.
- Validation happens in handler or request DTO layer, NOT in services.
- Validation errors must be user-friendly.

---

# 11. Module Rules

Each module MUST follow this structure:
internal/modules//
handler.go
service.go
repository.go (optional if needed)
routes.go
dto.go

Rules:
- handler → HTTP layer only
- service → business logic
- repository → DB abstraction (if needed)
- DTOs → request/response contracts

---

# 12. Security (OWASP-aligned)

- Always validate input
- Always enforce request body size limits
- Always use rate limiting on public endpoints
- Always sanitize user input before DB usage
- Never expose internal errors to client
- Use HTTPS in production
- Never log sensitive data

---

# 13. Performance Rules

- Avoid N+1 queries in Ent
- Use eager loading when needed
- Cache expensive read operations
- Do not block HTTP handlers with heavy computation
- Use background jobs for long-running tasks (if introduced later)

---

# 14. Code Style Rules

- Keep functions small and focused
- Prefer composition over inheritance
- Avoid deep nesting
- Use early returns
- Keep handlers thin
- No duplicated logic across modules

---

# 15. Testing Rules

- Services should be unit-testable via interfaces
- Use mocks for cache, DB, external services
- Avoid testing implementation details
- Focus on behavior

---

# 16. What NOT to do

- Do NOT introduce microservices
- Do NOT add global state
- Do NOT bypass container injection
- Do NOT put logic in handlers
- Do NOT couple Redis directly to business logic
- Do NOT return Ent models directly to API
- Do NOT skip validation
- Do NOT hardcode secrets

---

# 17. AI Agent Behavior Rules

When generating code:

- Always prefer clarity over cleverness
- Always follow existing project structure
- Always inject dependencies instead of creating them
- Always keep modules isolated
- Always assume production environment
- Never reduce security for convenience
- Always use existing abstractions first

---

# End of Rules
