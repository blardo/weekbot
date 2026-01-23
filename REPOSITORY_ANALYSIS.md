# WeekBot Repository Structure Analysis & Improvement Recommendations

## Current Structure Overview

```
weekbot/
├── .github/workflows/     # CI/CD
├── weekbot-go/            # Go Discord bot backend
│   ├── cmd/               # Application entry point
│   ├── internal/
│   │   ├── actions/       # Business logic actions
│   │   ├── api/           # HTTP API endpoints
│   │   ├── commands/     # Discord slash commands
│   │   ├── db/            # Database utilities (DUPLICATE)
│   │   ├── handlers/     # Discord event handlers
│   │   ├── models/       # Data models
│   │   ├── router/        # Handler/router configuration
│   │   └── services/     # External service integrations
│   └── weekbot-js/        # Next.js frontend (stub)
```

---

## Critical Issues

### 1. **Duplicate Database Packages** ⚠️ HIGH PRIORITY
- **Problem**: Two separate DB packages with overlapping functionality
  - `internal/db/db.go` - Has `GetDB()`, `GetCurrentSuggestions()`, `GetAllSuggestions()`
  - `internal/services/db.go` - Has `GetDB()` with reset functionality
- **Impact**: Confusion, potential bugs, inconsistent behavior
- **Recommendation**: 
  - Consolidate into `internal/services/db.go`
  - Move query functions to `internal/repositories/` or keep in models
  - Remove `internal/db/` package entirely

### 2. **Package Naming Inconsistency** ⚠️ MEDIUM
- **Problem**: `internal/actions/suggestion.go` contains `HandleWeekSuggestion` but was previously in `commands` package
- **Impact**: Confusing imports, unclear separation of concerns
- **Recommendation**: 
  - Keep `actions/` for business logic that doesn't directly handle Discord interactions
  - Keep `commands/` for Discord slash command handlers
  - Move `HandleWeekSuggestion` to `handlers/` since it's called from `ParseChatCommand`

### 3. **Missing Error Handling** ⚠️ MEDIUM
- **Problem**: Many functions use `println()` instead of proper logging, missing error returns
- **Impact**: Hard to debug, no structured logging
- **Recommendation**: 
  - Use structured logging (e.g., `logrus`, `zap`, or `slog`)
  - Replace all `println()` with proper log levels
  - Add error wrapping with context

### 4. **Global State Management** ⚠️ MEDIUM
- **Problem**: Multiple global singletons (`discordConnection`, `botInstances`, `globalTracker`)
- **Impact**: Hard to test, potential race conditions, unclear dependencies
- **Recommendation**: 
  - Use dependency injection pattern
  - Create a `Server` or `App` struct that holds all dependencies
  - Pass dependencies explicitly rather than using globals

---

## Structural Improvements

### 5. **Separation of Concerns**
**Current Issues:**
- Models contain business logic (e.g., `Poll.PerformRankedChoiceVoting()`)
- Handlers contain business logic
- No clear repository pattern

**Recommended Structure:**
```
internal/
├── cmd/              # CLI commands (if needed)
├── config/           # Configuration management
├── models/           # Data models only (structs, no business logic)
├── repositories/     # Data access layer (DB queries)
├── services/         # Business logic services
│   ├── poll/         # Poll-related business logic
│   ├── suggestion/   # Suggestion-related business logic
│   └── voting/       # Voting algorithm logic
├── handlers/         # Discord event handlers (thin layer)
├── commands/         # Discord slash command handlers
├── api/              # HTTP API handlers
└── discord/          # Discord client wrapper
```

### 6. **Constants and Configuration**
**Current Issues:**
- Magic numbers scattered (e.g., `>= 3`, `>= 5`)
- Hard-coded strings ("week-name", "bd", acceptable weeks)
- Test thresholds mixed with production code

**Recommendation:**
```go
// internal/config/constants.go
package config

const (
    MinSuggestionsToStartPoll = 3
    MinUpdicksToQualify      = 3
    MinBallotsToEndPoll      = 5
    
    WeekNameChannelName      = "week-name"
    QualifyingEmoji         = "bd"
    ConfirmationEmoji       = "👍"
)

var AcceptableWeekSuffixes = []string{"week", "week.", "week!", "week?"}
```

### 7. **Testing Infrastructure**
**Current Issues:**
- Only one test file exists
- No test utilities or mocks
- Hard to test due to global state

**Recommendation:**
- Add `internal/testing/` package with test utilities
- Create interfaces for Discord service (enables mocking)
- Add integration tests for critical flows
- Use table-driven tests for voting algorithm

### 8. **Error Handling Strategy**
**Current Issues:**
- Inconsistent error handling
- Some functions panic, others return errors
- No error wrapping or context

**Recommendation:**
- Define custom error types
- Use `fmt.Errorf()` with `%w` for wrapping
- Create error constants for common errors
- Add error context (guild ID, user ID, etc.)

---

## Code Quality Improvements

### 9. **Logging Standardization**
**Current Issues:**
- Mix of `fmt.Println()`, `println()`, and `log.Printf()`
- No log levels
- No structured logging

**Recommendation:**
```go
// Use structured logging
logger.WithFields(log.Fields{
    "guild_id": guildID,
    "user_id": userID,
    "message_id": messageID,
}).Info("Processing suggestion")
```

### 10. **Type Safety**
**Current Issues:**
- String-based channel lookups
- Magic strings for custom IDs
- No validation of user input

**Recommendation:**
- Create types for Discord IDs
- Use constants for custom IDs
- Add input validation

### 11. **Database Migrations**
**Current Issues:**
- Auto-migration only
- No migration history
- No rollback capability

**Recommendation:**
- Use a migration tool (e.g., `golang-migrate`)
- Version migrations
- Add migration tests

---

## Project Organization

### 12. **Documentation**
**Missing:**
- README.md for weekbot-go
- API documentation
- Architecture decision records
- Contributing guidelines

**Recommendation:**
- Add comprehensive README
- Document environment variables
- Add setup instructions
- Document deployment process

### 13. **Development Tools**
**Missing:**
- Makefile for common tasks
- Pre-commit hooks
- Linting configuration (golangci-lint)
- Code formatting enforcement

**Recommendation:**
```makefile
# Makefile
.PHONY: run test lint format build

run:
	go run ./cmd/main.go

test:
	go test ./...

lint:
	golangci-lint run

format:
	go fmt ./...
```

### 14. **Environment Management**
**Current Issues:**
- Secrets in root directory (`dtk`, `appid`)
- No `.env.example`
- No validation of required env vars

**Recommendation:**
- Add `.env.example` with placeholders
- Validate required env vars on startup
- Use `viper` or similar for config management
- Document all required variables

### 15. **Git Hygiene**
**Issues:**
- `.DS_Store` files committed
- Secrets potentially in git history
- No `.gitattributes`

**Recommendation:**
- Add `.DS_Store` to `.gitignore`
- Add `.gitattributes` for line endings
- Consider using `git-secrets` or `truffleHog` to scan history

---

## Performance & Scalability

### 16. **Database Connection Management**
**Current Issues:**
- New DB connection per query in `internal/db/`
- No connection pooling
- No connection lifecycle management

**Recommendation:**
- Reuse connections from Bot struct
- Implement connection pooling
- Add connection health checks

### 17. **Caching Strategy**
**Current:**
- Channel cache implemented (good!)
- Message deduplication cache (good!)

**Recommendation:**
- Consider caching poll state
- Cache suggestion lookups
- Add cache invalidation strategy

### 18. **Concurrency Safety**
**Current:**
- Some mutex usage (good)
- Global state without proper locking in some places

**Recommendation:**
- Audit all shared state
- Use `go vet` and race detector
- Document thread-safety guarantees

---

## Deployment & Operations

### 19. **Health Checks**
**Missing:**
- No health check endpoint
- No readiness probe
- No graceful shutdown for DB connections

**Recommendation:**
- Add `/health` and `/ready` endpoints
- Implement graceful shutdown
- Add shutdown hooks for cleanup

### 20. **Monitoring & Observability**
**Missing:**
- No metrics collection
- No distributed tracing
- Limited error tracking

**Recommendation:**
- Add Prometheus metrics
- Add structured logging for events
- Consider error tracking (Sentry, etc.)

---

## Priority Action Items

### Immediate (This Sprint)
1. ✅ Consolidate duplicate DB packages
2. ✅ Fix package naming inconsistencies
3. ✅ Add proper logging infrastructure
4. ✅ Remove `.DS_Store` files
5. ✅ Create constants file

### Short Term (Next Sprint)
6. Refactor to dependency injection
7. Separate business logic from handlers
8. Add comprehensive error handling
9. Create repository pattern
10. Add Makefile and dev tools

### Medium Term (Next Month)
11. Add comprehensive tests
12. Implement proper migrations
13. Add monitoring/metrics
14. Improve documentation
15. Add health checks

---

## Recommended File Structure (Refactored)

```
weekbot-go/
├── cmd/
│   └── main.go
├── internal/
│   ├── config/
│   │   ├── constants.go
│   │   └── config.go
│   ├── models/
│   │   ├── ballot.go
│   │   ├── bot.go
│   │   ├── poll.go
│   │   ├── suggestion.go
│   │   └── voter.go
│   ├── repositories/
│   │   ├── ballot.go
│   │   ├── poll.go
│   │   └── suggestion.go
│   ├── services/
│   │   ├── db/
│   │   │   └── db.go
│   │   ├── discord/
│   │   │   └── discord.go
│   │   ├── poll/
│   │   │   └── service.go
│   │   └── voting/
│   │       └── ranked_choice.go
│   ├── handlers/
│   │   ├── handlers.go
│   │   ├── deduplication.go
│   │   └── channel_cache.go
│   ├── commands/
│   │   ├── poll.go
│   │   └── ping.go
│   ├── actions/
│   │   └── suggestion.go
│   ├── api/
│   │   └── ping.go
│   └── router/
│       └── router.go
├── migrations/
│   └── (migration files)
├── Makefile
├── .env.example
├── README.md
└── go.mod
```

---

## Conclusion

The codebase has a solid foundation but needs structural improvements for maintainability, testability, and scalability. The most critical issues are duplicate packages and global state management. Addressing these will make the codebase more robust and easier to work with.
