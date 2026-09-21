SB1-01 
  - Initialize module slotbook. The entry point lives in cmd/slotbook-api; every other package lives
    under internal/. No package may be named utils, common or helpers.


SB1-02
  - SB1-02 Provide a Makefile with targets run, test, lint and fmt. make test MUST run all tests with the
    race detector enabled.


SB1-03
  - The listen address comes from env var SLOTBOOK_ADDR (default :8080). Only main reads the
    environment; packages receive values through constructors.

SB1-04 
  - GET /healthz returns 200 with JSON fields status and version; the version is injected at build
    time via linker flags and defaults to dev.

SB1-05  
  - A Resource has id (UUID), name (1–80 chars, unique case-insensitively), kind (room, desk or
    studio), capacity (1–500), created_at and updated_at (RFC 3339, UTC).

SB1-06
  - Implement POST /v1/resources (201), GET /v1/resources (200), GET /v1/resources/{id}
    (200), PATCH /v1/resources/{id} (200) and DELETE /v1/resources/{id} (204).
