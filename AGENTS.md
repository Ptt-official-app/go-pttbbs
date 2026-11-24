# Repository Guidelines

## Project Overview
`github.com/Ptt-official-app/go-pttbbs`: a Go implementation of [ptt/pttbbs](https://github.com/ptt/pttbbs) exposing an HTTP/JSON API (Gin) on top of the legacy PTT BBS on-disk data (`.PASSWDS`, `.BRD`, `boards/*/.DIR`, `home/`) and its SysV shared memory (SHM) and semaphores. It is the backend of [pttbbs-backend](https://github.com/Ptt-official-app/pttbbs-backend). Go 1.25 (`go.mod`). Also contains gRPC stubs for `boardd` and `mand` (generated `*.pb.go`, do not hand-edit).

## Architecture & Data Flow
Client -> Gin route (`initgin/init_gin.go`, prefix `/v1`) -> `api.XWrapper` -> binder (`api.JSON/Query`, `api.LoginRequiredJSON/Query/PathJSON/PathQuery`) -> `api.X` -> `bbs.X` -> `ptt.X` -> `cache` (SHM/in-memory maps), `cmbbs` / `cmsys` (passwd, record files, locks), `ptttype` / `types` (raw C-compatible structs, Cstr, Big5).
- `api.process*` reads `Host` and `X-Forwarded-For` (the latter must be non-empty) and calls `processResult`, which maps errors to HTTP 400/401/403 (default 500).
- Login-required binders derive the user from the JWT in the `Authorization` header; a missing/invalid token falls back to `guest`. Several list/read endpoints accept `system=true` to act as SYSOP: the service is meant to be reached only by a trusted middleware.
- `bbs` converts API-facing types (`BBoardID`, `UUserID`, `ArticleID`, summaries) to/from `ptttype` raw types; `ptt` ports logic from the C `mbbsd`. Comments reference the original C source URLs.
- `types.IS_ALL_GUEST` selects `initGinAllGuest()` (read-only, no SHM user state) or `initGin()` (full API).
- Startup (`go-pttbbs/main.go`, `main_init.go`): `initgin.InitAllConfig(ini)` -> `cache.Init(SHM_KEY, ...)` -> (non all-guest) `cache.LoadUHash`, `cache.ReloadBCache` when `IS_NEW_SHM`, `cache.AttachCheckSHM` -> `router.Run(HTTP_HOST)`.

## Key Directories
- `api/` one file per endpoint (`X_R` route const, `XParams`, `XPath`, `XResult`, `XWrapper`, `X`), plus `api.go`, `api_login_required*.go`, `auth_utils.go` (JWT), `errors.go`
- `bbs/` service layer between `api` and `ptt`; `ptt/` ported BBS logic (`bbs.go`, `board_list.go`, `article_list.go`, `register.go`, `fav/`)
- `cache/` SHM and in-memory cache (boards, users, utmp); `shm/` and `sem/` cgo wrappers for SysV shm/sem; `cmbbs/` passwd and file helpers (`import "C"`), `cmsys/` records, locks, strings; `crypt/` pure-Go fcrypt
- `ptttype/` C-compatible raw structs, constants, config vars (`00-config*.go`); `types/` Cstr, Big5, time, file helpers, `ansi/`, `proto/` (generated rune protobuf)
- `boardd/`, `mand/` gRPC clients/mocks (proto-generated); `fromd/`; `go-pttbbs/` the server `main` package; `initgin/` router and config init; `configutil/` viper helpers
- `c-pttbbs/` C helper that prints SHM struct offsets (`shm_offset.c`); `pttbbs/` git submodule (ptt/pttbbs), do not edit
- `docs/` config templates (`docs/config`), sample BBS data (`docs/examples`, `docs/etc`); `docker/`, `scripts/`, `apidoc/` (Python swagger generator, `proto/rune.proto`)

## Development Commands
```
go build ./...                       # build everything
./scripts/test.sh                    # go build ./... && go test -p 1 -v ./... -cover, then ipcrm test shm/sem
go test -p 1 ./bbs ./ptt/...         # one package; -p 1 is required (shared SHM/sem keys, shared testcase files)
./scripts/coverage.sh                # coverage report (needs gotest)
golangci-lint run                    # same as scripts/lint.sh and CI (golangci-lint v2.7)
cd docs; tar -zxvf bbs-2025-03-22.tar.gz; cd ..; ./scripts/run-dev.sh   # build with -tags dev and run on :3456
cd go-pttbbs && go build -tags docker  # production-like docker build (see docker/Dockerfile)
./scripts/protoc.sh | protoc_board.sh | protoc_man.sh                   # regenerate protobuf stubs
./scripts/swagger.sh                 # swagger UI from apidoc/ (needs python virtualenv)
```
Version is injected with `-ldflags -X .../types.GIT_VERSION=... -X .../types.VERSION=...`.

## Code Conventions & Common Patterns
- **Endpoint pattern** (copy `api/create_article.go`): `const X_R = "/board/:bid/..."`; `XParams` (tags `json`/`form`/`url`), `XPath` (tag `uri`), `XResult`; `func XWrapper(c *gin.Context)` calls a binder; `func X(remoteAddr string, uuserID bbs.UUserID, params interface{}, path interface{}, c *gin.Context) (result interface{}, err error)` type-asserts params/path and returns `ErrInvalidParams` / `ErrInvalidPath` on mismatch. Register in `initgin/init_gin.go` (and the all-guest router if it is read-only).
- **Errors**: sentinel `errors.New` in each package's `errors.go`; new errors that need a non-500 status must be added to `api.processResult`.
- **Config**: three files per module: `00-config.go` (vars), `config.go` (read from ini via viper, `setXConfig("KEY", KEY)`), `config_util.go` (helpers). ini template: `docs/config/*.ini`. Compile-time consts for `ptttype` (IDLEN, PASSLEN, ...) live in `ptttype/00-config-*.go` selected by build tags (`dev`, `production`, `docker`; default `!dev && !production && !docker`); see README "Config".
- **Raw types**: fixed-size byte arrays (`UserID_t`, `Filename_t`, `Title_t`) mimic C strings; use `types.Cstr*` helpers, and check `.DIR` record indexes (`SortIdx` is 1-based, `SortIdxInStore` is 0-based).
- Naming: ALL_CAPS for config vars, constants, routes; snake_case file names.
- Formatting: gofmt/gofumpt (README); lint config in `.golangci.yaml` (gosec, misspell US, revive without var-naming, nolintlint); tests are not linted.

## Important Files
`go-pttbbs/main.go`, `go-pttbbs/main_init.go`, `initgin/init_gin.go`, `api/api.go`, `api/api_utils.go` (`processResult`), `api/auth_utils.go`, `ptttype/00-config.go`, `cache/shm.go`, `cmsys/record.go` (.DIR record access), `.golangci.yaml`, `docs/config/`, `docker/Dockerfile`, `.github/workflows/go.yml`.

## Runtime/Tooling Preferences
- Go per `go.mod`; `docker/Dockerfile` builds with `golang:1.25.x` and needs a C toolchain (cgo is used by `shm` and `sem`).
- Runtime needs access to BBSHOME data and SysV SHM/semaphore (`SHM_KEY`, `PASSWDSEM_KEY` in ptttype config). No MongoDB or Redis.
- `scripts/test.sh` and CI clean leftover test SHM/sem with `ipcrm -S 0x00007ffb` / `ipcrm -M 0x0000fffe`.

## Testing & QA
- testing + testify; table-driven tests (`TestXxx`) in `<source>_test.go` beside the source; `testutil.TDeepEqual` for deep comparison.
- Each package that touches shared state has `testinit.go` (`SetIsTest/UnsetIsTest`, mutex-guarded, test keys/paths) and `testinit_test.go` with `setupTest()` / `teardownTest()`; tests start with `setupTest(); defer teardownTest()`. Test data lives in `<pkg>/testcase/` (e.g. `cache/testcase/.PASSWDS1`, `.BRD1`, `boards/`); setup copies `.PASSWDS1`/`.BRD1` to the working `.PASSWDS`/`.BRD` and teardown removes them.
- Always run with `-p 1`. CI (`.github/workflows/go.yml`) runs golangci-lint, `go build`, then `go test -p 1 -v ./... -cover` per package directory with `ipcrm` cleanup in between, and a final coverage run uploaded to codecov. No external services are required.
- Codecov (`codecov.yml`): project target 50% (5% threshold); ignores `apidoc`, `scripts`, `pttbbs`, `swagger`, `testutil`.
