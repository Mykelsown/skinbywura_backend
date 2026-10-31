# syntax=docker/dockerfile:1

# ---- build stage: compiles a single static binary ----
FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api \
 && CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/importproducts ./cmd/importproducts

# ---- run stage: tiny image, non-root user ----
FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata \
 && adduser -D -u 10001 app
WORKDIR /app
COPY --from=build /out/api /app/api
COPY --from=build /out/importproducts /app/importproducts
USER app
ENV APP_ENV=production PORT=8080
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s --start-period=20s CMD wget -qO- http://127.0.0.1:${PORT}/api/health || exit 1
CMD ["/app/api"]
