# go-logger

Shared Go module — structured JSON logger used by every Uniplaces Go service. Wraps a forked
`logrus` (`github.com/uniplaces/logrus`) behind a small, opinionated surface: a singleton
instance, a fluent per-log field builder, and a process-wide "default fields" list that gets
merged into every entry.

Consumed via `import logger "github.com/uniplaces/go-logger"`. Module path
`github.com/uniplaces/go-logger`, package name `go_logger`.

## Who consumes this

All 10 Go services pin `v0.7.0` (the current latest tag) in `go.mod`:

| Service | Typical use |
|---|---|
| `aggregate-graphql`, `offer-aggregator`, `authentication`, `search-api`, `email-gateway`, `communications` | HTTP request middleware + recovery middleware + infrastructure-adapter error logging |
| `search-ingestors` | CLI entry log, per-message handler logs, bulk ETL progress |
| `salesforce-integration`, `integrations-enricher-backend`, `offer-aggregate-cdn-version-updater` | CLI / SQS-worker logging |

`v0.7.0` (commit `35fcd72`, Jan 2020) is the latest release.

## Public surface

```
logger.Init(logger.NewConfig(env, level)) error      // singleton init; call once at boot
logger.InitWithInstance(custom Logger) error         // test seam — inject a Logger impl

logger.AddDefaultField(key, value, isContextField)   // appended to every subsequent log
logger.Builder() builder                             // start a per-log field chain

// Module-level shortcuts (use Builder() under the hood):
logger.Error(err)            // err-level, captures stack trace
logger.Warning(msg)
logger.Info(msg)
logger.Debug(msg)

// Builder terminals:
b := logger.Builder().AddField("k", v).AddContextField("ck", cv)
b.Error(err) / b.Warning(msg) / b.Info(msg) / b.Debug(msg)

// Logger interface (interface.go) — what InitWithInstance accepts:
type Logger interface {
    ErrorWithFields(err error, fields map[string]interface{})
    WarningWithFields(message string, fields map[string]interface{})
    InfoWithFields(message string, fields map[string]interface{})
    DebugWithFields(message string, fields map[string]interface{})
}
```

## File map

| File | Role |
|---|---|
| `logger.go` | Singleton (`instance Logger`, `sync.Once`), module-level `Init`/`Error`/`Warning`/`Info`/`Debug`, `AddDefaultField` |
| `builder.go` | `builder` struct + fluent `AddField` / `AddFields` / `AddContextField`; merges default fields on each terminal call |
| `interface.go` | `Logger` interface (4 `*WithFields` methods) |
| `config.go` | `Config{environment, level}` + `NewConfig` |
| `internal/logrus.go` | Concrete `logrusLogger` implementation — JSON formatter, stack-trace extraction (`pkg/errors.stackTracer` aware + runtime fallback), level mapping |
| `internal/logrus_test.go` | Logrus binding tests |
| `builder_test.go` / `logger_test.go` / `config_test.go` | Behaviour tests for the public surface |
| `Gopkg.toml` / `Gopkg.lock` | Legacy `dep` manifests (see "Dependency tooling" below) |

## Mandatory default fields

`Init` (and `InitWithInstance`) calls `addMandatoryDefaultFields()`, which appends these to
`defaultFields` so every log entry carries them:

| Field | Source | Example |
|---|---|---|
| `type` | constant `"app"` | `"app"` |
| `app-id` | `$APPID` env var | `"aggregate-graphql"` |
| `env` | `$GOENV` env var | `"staging"` |
| `git-hash` | `$GITHASH` env var | `"a1b2c3d"` |

If those env vars are empty, the field is still emitted with an empty value (matches the
example log in the operations dashboards where `git-hash` shows blank).

## Default fields vs builder fields

| Mechanism | Lifetime | API | Use it for |
|---|---|---|---|
| **Mandatory defaults** | Set at `Init` from env vars | (private) | `type`, `app-id`, `env`, `git-hash` |
| **User-added defaults** | Set anywhere at boot via `AddDefaultField` | `logger.AddDefaultField(k, v, isContextField)` | Static process-wide values (e.g. `ingestion_run_type` in search-ingestors) |
| **Builder fields** | Per-log line only | `logger.Builder().AddField(k, v).AddContextField(ck, cv).Info(msg)` | Anything per-request, per-message, per-iteration |

`AddDefaultField` mutates a package-level slice. Every service today calls it during DI
setup, before any concurrent logging starts.

The "context" key is a nested map: any field passed via `AddContextField` (or
`AddDefaultField(..., isContextField=true)`) lands under `"context": { ... }` in the JSON
output instead of as a top-level key. `README.md` recommends it for "contextual information
when logging an error", e.g. the entity ID that failed to load.

## Stack-trace handling (`internal/logrus.go`)

- **Error level always captures a stack trace.** `stackTraceLevels` map controls this; only
  `ErrorLevel` is `true` today.
- If the error implements `pkg/errors.stackTracer` (`StackTrace() errors.StackTrace`), the
  logger uses *that* trace, walking the `Cause()` chain via `firstStackTracerInErrorChain` to
  find the deepest wrapped stack-tracing error.
- Otherwise it falls back to a `runtime.Caller` walk, skipping frames matching this list:
  `github.com/uniplaces/go-logger`, `github.com/gin-gonic`, `autogenerated`,
  `go/src/runtime/asm_amd64.s`, `go/src/net/http/server.go`.
- Trace lines are formatted `%s:%d` (runtime) or `%+v` (pkg/errors frames) and emitted under
  the `stack_trace` JSON key as an array of strings.

Warning/Info/Debug entries do not emit a stack trace — `stackTraceLevels` sets them all to
`false`. Flip a value to enable; nothing else needs to change.

## Singleton model and test seam

`Init` is enforced once: a second call returns `"logger cannot be initialized more than once"`
and is a no-op. The `sync.Once` inside the function is belt-and-braces.

`InitWithInstance(custom Logger)` exists to inject a test or in-memory implementation. Tests
in this repo use it with `internal.NewLogrusLogger(level, env, &bytes.Buffer{}, ...)` and a
`resetInstance()` helper. Consumers do the same pattern.

Boot-time wiring is the only valid mutation point for `instance`.

## Dependency tooling — `dep`

This repo predates Go modules. The dep manifests are checked in:

```
Gopkg.toml   constraints (logrus v1.1.0, testify 1.2.0, pkg/errors 0.8.0)
Gopkg.lock   resolved revisions
vendor/      .gitignored
```

There is no `go.mod`. Go-modules consumers resolve via the `v0.7.0` tag — Go's module proxy
synthesises a pseudo-module for any tagged repo without a `go.mod`.

## Forked logrus (`github.com/uniplaces/logrus`)

We pin a Uniplaces fork of logrus at `v1.1.0`. The fork's reason-for-being is recorded in the
historical PRs (`Update to last commit (forked logrus)`, `Use forked logrus version`) and is
RFC3164-compliance related. Log shipping (Graylog) depends on the formatter behaviour.

## Versioning

Tags are sequential semver-ish. Latest is `v0.7.0`. Consumer pinning is done by `go.mod`
require line, e.g. `github.com/uniplaces/go-logger v0.7.0`.

**Backward compatibility:** field names in JSON output are part of the contract — log
shippers, dashboards, and alerting rules grep on `app-id`, `env`, `git-hash`, `type`,
`stack_trace`, `context`. Coordinating downstream changes is required before renaming any of
those.

## Repo conventions

- **Package name:** `go_logger` (underscore). The repo dir is `go-logger` (hyphen). The repo
  dir / package name mismatch is historical.
- **Tests:** Same-package tests (`package go_logger`, `package internal`). `t.Parallel()` is
  used in `config_test.go` and `internal/logrus_test.go` but **not** in `logger_test.go`
  because those tests share `instance` (global mutable state).
- **Trailing blank line:** every file ends with a blank line.
- **No build tags, no Makefile, no Dockerfile.** This repo only ships source; tooling lives
  in the consuming service.