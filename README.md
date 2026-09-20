## Slotbook


A System that lets organizations publish bookable resources (meeting rooms, desks, studios) and lets their
members hold, confirm and cancel time slots. It starts as an in-memory CRUD service and grows into a
production-shaped API: PostgreSQL with database-enforced overlap rules, JWT authentication with
refresh-token rotation, role-based authorization, idempotent writes, optimistic locking, rate limiting, a
transactional outbox, Redis caching, full observability
