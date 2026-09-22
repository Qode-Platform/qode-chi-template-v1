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

## Rule: everything under BASE_PATH

The fleet serves this app behind a proxy at `BASE_PATH=/direct/<agent>:<port>`,
and that prefix is forwarded **unchanged** - nginx does not strip it. So the app
sees the full prefixed path on every request, and everything it emits must carry
the prefix too.

- **Register every route on the group returned from `basePath()`**, never on the
  root router. In this template that group is ``root.Mount(basePath(), app)` in `newRouter()``; a route added
  straight to the root router answers only when `BASE_PATH` is empty and 404s in
  the fleet.
- **Route patterns are group-relative.** Write `"/health"`, not
  `"/direct/agent-7:3000/health"` - the group supplies the prefix.
- **Build every absolute URL you hand a client by prefixing `basePath()`.** That
  covers redirect targets (`Location:`), links and form actions in any HTML you
  render, and asset paths (CSS, JS, images). A bare `"/login"` or `"/static/app.css"`
  escapes the prefix and lands outside the app.
- Relative URLs and prefixed ones are both fine; a leading-slash literal is the
  thing to look for in review.
- `basePath()` returns `""` when `BASE_PATH` is unset, so the same code serves
  standalone at the host root.
- **chi survives the colon.** chi uses `{param}` syntax, so the `:` in
  `/direct/<agent>:<port>` is a literal path character and mounts cleanly.
  path-to-regexp-based routers (Express, and the Node/TS servers built on it)
  read `:3000` as a route parameter and crash at boot on the same prefix.
