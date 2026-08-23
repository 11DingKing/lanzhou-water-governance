FROM golang:1.23-bookworm AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 go build -o /out/lanzhou ./cmd/server

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=build /out/lanzhou /app/lanzhou
COPY migrations /app/migrations
ENV ADDR=:8080 DB_PATH=/data/lanzhou.db
RUN mkdir -p /data
EXPOSE 8080
ENTRYPOINT ["/app/lanzhou"]
