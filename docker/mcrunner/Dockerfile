# ---- Build stage ----
FROM golang:1.25-alpine AS builder

ARG GIT_COMMIT
ARG GIT_DATE
ARG GIT_TAG

WORKDIR /workdir

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build -ldflags "\
    -w -s \
    -X 'main.gitCommit=${GIT_COMMIT}' \
    -X 'main.gitDate=${GIT_DATE}' \
    -X 'main.gitTag=${GIT_TAG}'" \
  -o mcrunner ./cmd/mcrunner

# ---- Runtime stage ----
FROM ubuntu:24.04

RUN apt-get update && apt-get install -y \
  ca-certificates \
  tzdata && \
  update-ca-certificates \
  default-jre

WORKDIR /workdir

COPY --from=builder /workdir/mcrunner /usr/bin/mcrunner

EXPOSE 3000 25565

ENTRYPOINT ["/usr/bin/mcrunner"]
