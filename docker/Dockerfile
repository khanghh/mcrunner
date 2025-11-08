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
  -o mcrunner main.go

# ---- Runtime stage ----
FROM ubuntu:24.04

RUN apt-get update && apt-get install -y \
  ca-certificates \
  tzdata \
  cron \
  default-jre && \
  update-ca-certificates

WORKDIR /minecraft

COPY --from=builder /workdir/mcrunner /usr/bin/mcrunner

COPY ./docker/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 3000 25565

ENTRYPOINT ["/entrypoint.sh"]
