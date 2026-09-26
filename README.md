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
| start | `env PORT="$PORT" ./.bin/app` |

Listens on `$PORT` and serves at the root of its own hostname
(`https://<hash>.<FLEET_APP_DOMAIN>/`); health check hits `/health`.
