# Built by .github/workflows/deploy.yml and pushed to Artifact Registry.
# Adapted from the fleet's go stack pack.
#
# Deviations, and why:
#   - golang:1.25-alpine, not the pack's 1.23: `go mod tidy` resolves this
#     module to go 1.25 (gin 1.12 alone requires >= 1.25).
#
# BASE_PATH is not baked in - the binary reads it from the environment.
FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/app ./cmd/app

FROM gcr.io/distroless/static-debian12:nonroot AS runtime
WORKDIR /app
ARG BUILD_ID=""
ENV PORT=8080 BUILD_ID=$BUILD_ID
COPY --from=build /out/app /app/app
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app/app"]
