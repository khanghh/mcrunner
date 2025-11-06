# ---- Build go app----
FROM golang:1.25-alpine AS go-builder

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
  -o manager ./cmd/manager

# ---- Build vscode web ----
FROM mcr.microsoft.com/devcontainers/typescript-node:16-bullseye AS vscode-builder

RUN apt-get update && apt-get install -y libsecret-1-dev libxkbfile-dev
RUN git config --system codespaces-theme.hide-status 1
RUN yarn global add node-gyp@9.3.1 

USER node

WORKDIR /workdir

COPY --chown=node:node ./web/vscode-web .
RUN ls -lah /workdir

RUN ./build.sh

# ---- Build manager web----
FROM node:20-bullseye AS web-builder

ENV CI=true

RUN npm install --global corepack@latest

WORKDIR /workdir

COPY --chown=node:node ./web .
RUN rm -rf ./src/proto
COPY --chown=node:node ./internal/websocket/proto ./src/proto

RUN pnpm install && pnpm run build

# ---- Runtime stage ----
FROM scratch

WORKDIR /app

COPY --from=go-builder /workdir/manager /app/manager
COPY --from=web-builder /workdir/dist /app/dist
COPY --from=vscode-builder /workdir/dist /app/dist/vscode

ENTRYPOINT ["/app/manager"]
