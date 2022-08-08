# Balance Backend

[![License](https://img.shields.io/github/license/yyq1025/balance-backend)](https://github.com/yyq1025/balance-backend/blob/main/LICENSE)
[![Lines of code](https://img.shields.io/tokei/lines/github/yyq1025/balance-backend)](https://github.com/yyq1025/balance-backend)
[![Go Reference](https://pkg.go.dev/badge/github.com/yyq1025/balance-backend.svg)](https://pkg.go.dev/github.com/yyq1025/balance-backend)
[![Go Report Card](https://goreportcard.com/badge/github.com/yyq1025/balance-backend)](https://goreportcard.com/report/github.com/yyq1025/balance-backend)
[![Go CI](https://github.com/yyq1025/balance-backend/actions/workflows/ci.yml/badge.svg)](https://github.com/yyq1025/balance-backend/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/yyq1025/balance-backend/branch/main/graph/badge.svg?token=XHSJ1TK3KQ)](https://codecov.io/gh/yyq1025/balance-backend)
[![Renovate status](https://img.shields.io/badge/renovate-enabled-brightgreen?logo=renovatebot)](https://github.com/yyq1025/balance-backend/issues/12)

## Environment Variables

| Name | Required | Default | Description |
| ---- | -------- | ------- | ----------- |
| `AUTH0_DOMAIN` | Yes || Auth0 domain, learn more [here](https://auth0.com/docs/quickstart/backend/golang) |
| `AUTH0_AUDIENCE` | Yes || Auth0 audience, learn more [here](https://auth0.com/docs/quickstart/backend/golang) |
| `PORT` | No | `8080` | Http port to serve |
| `TIMEOUT` | No | `3.5s` | Response timeout |
| `DB_HOST` | No | `localhost` | PostgreSQL host |
| `DB_PORT` | No | `5432` | PostgreSQL port |
| `DB_USER` | No | `postgres` | PostgreSQL username |
| `DB_PASSWORD` | No | `postgres` | PostgreSQL password |
| `DB_NAME` | No | `postgres` | PostgreSQL database name |
| `REDIS_HOST` | No | `localhost` | Redis host |
| `REDIS_PORT` | No | `6379` | Redis port |
| `REDIS_USER` | No | `""` | Redis username |
| `REDIS_PASSWORD` | No | `""` | Redis password |
