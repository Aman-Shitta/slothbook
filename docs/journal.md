SB1-01 
  - Initialize module slotbook. The entry point lives in cmd/slotbook-api; every other package lives
    under internal/. No package may be named utils, common or helpers.


SB1-02
  - SB1-02 Provide a Makefile with targets run, test, lint and fmt. make test MUST run all tests with the
    race detector enabled.
