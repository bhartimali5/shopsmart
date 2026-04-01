# ShopSmart Test Strategy

## Scope
Covers unit, integration, and E2E testing for: cart, order, payment, user, and notification services.
Excludes infrastructure (RabbitMQ, SQLite setup) and third-party dependencies.

---

## Risk Areas

| Area | Risk | Reason |
|---|---|---|
| Cart total calculation | High | Directly affects order amount |
| Place order flow | High | Multi-step, involves DB + RabbitMQ |
| Payment processing | High | Financial outcome |
| JWT validation | High | Security boundary |
| Order status transitions | High | Data consistency across services |
| User signup/login | Medium | Standard auth, well-understood |
| Email notification | Low | Non-critical, fire-and-forget |
| List products/categories | Low | Read-only, no side effects |

---

## Test Types Per Layer

### Unit (15–20 tests)
Focus: pure logic, no external dependencies.
- Cart: total price calculation, item quantity updates
- Payment: success/failure random logic
- JWT: token generation and verification
- Input validation: required fields, formats

### Integration (8–10 tests)
Focus: service + database interaction, middleware behaviour.
- Cart → SQLite: add/remove items persisted correctly
- Order → SQLite: order saved with correct status
- Payment consumer → SQLite: payment record created after event
- Internal auth middleware: valid key passes, invalid key returns 401
- User service: fetch email by user ID returns correct data

### E2E (2–3 tests)
Focus: full user journey across services.
- Happy path: signup → login → add to cart → place order → payment succeeds → order status updated
- Payment failure path: order placed → payment fails → order status reflects failure
- Auth failure: request without token is rejected at protected endpoints

---

## What Won't Be Tested & Why

| Area | Reason |
|---|---|
| RabbitMQ broker itself | Infrastructure concern, not application logic |
| SQLite engine | Third-party, assumed reliable |
| Swagger docs | Generated code, not business logic |
| Mock email sender (log output) | Placeholder, no logic to verify |
| OS/environment setup | Out of scope for application tests |

---

## Test Data Strategy

### Unit Tests — Inline Literals
Data is defined directly in the test function. No external setup needed.
```go
cart := Cart{Items: []CartItem{{Price: 100.0, Quantity: 2}}}
assert.Equal(t, 200.0, cart.TotalPrice())
```

### Integration Tests — DB Fixtures / Factories
Use a separate test SQLite DB (`:memory:` or `test.db`).
Seed required records before each test, clean up after.
```go
// factory helper
func createTestUser(db *sql.DB) models.User { ... }
func createTestCart(db *sql.DB, userID string) models.Cart { ... }
```
Each test is self-contained — no shared state between tests.

### E2E Tests — Seeded Test Users
A fixed set of test users pre-loaded before the E2E suite runs.
```
test_user@shopsmart.com / testpass123   (role: user)
test_admin@shopsmart.com / adminpass123 (role: admin)
```
Seeding runs once before the suite, teardown runs after.
E2E tests use these credentials to simulate real user flows across services.

---

## Summary

```
        /\        E2E         2–3 tests   (full journey)
       /  \
      /----\      Integration 8–10 tests  (service + DB + middleware)
     /      \
    /--------\    Unit        15–20 tests (pure logic)
```

Pyramid shape confirmed. High-risk areas covered at unit + integration level.
Low-risk areas (notifications, read-only endpoints) covered minimally or not at all.
