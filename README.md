# chi template

Provisioned from [`Qode-Platform/fleet-template-v1`](https://github.com/Qode-Platform/fleet-template-v1) - the fleet
lifecycle contract with a chi starter on top.

## Verified

Built and tested locally on Go 1.23.4 (toolchain auto-upgraded to 1.25):
`go build ./...` and `go test ./...` both pass.

## Fleet lifecycle

| step | command |
|---|---|
| install | `go mod download` |
| build | `go build -o ./.bin/app ./cmd/app` |
| start | `env PORT="$PORT" BASE_PATH="$BASE_PATH" ./.bin/app` |

Listens on `$PORT`; health check hits `/health`, and every route is mounted
under `$BASE_PATH` when the fleet injects one.

## Notes

- chi uses {param} syntax, so the colon in BASE_PATH is a literal path character - the same prefix crashes path-to-regexp routers.
