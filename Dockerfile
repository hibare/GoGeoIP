ARG GOLANG_VERSION=1.26.0

FROM golang:${GOLANG_VERSION}-alpine AS base

# ================== Build App ================== #
FROM base AS build

# prefer ash on Alpine so -o pipefail behaves as expected
SHELL ["/bin/ash", "-o", "pipefail", "-c"]

ARG VERSION=unknown
ARG BUILD_TIMESTAMP=unknown
ARG COMMIT_HASH=unknown

# Install healthcheck cmd
# hadolint ignore=DL3018
RUN apk update \
    && apk add --no-cache ca-certificates curl cosign \
    && rm -rf /var/cache/apk/* \
    && curl -sfL https://raw.githubusercontent.com/hibare/go-docker-healthcheck/main/install.sh | sh -s -- -d -v -b /usr/local/bin

WORKDIR /src/

COPY . /src/

# hadolint ignore=DL3018
RUN apk --no-cache add ca-certificates \
    && CGO_ENABLED=0 go build -ldflags "-X github.com/hibare/Waypoint/internal/constants.Version=$VERSION -X github.com/hibare/Waypoint/internal/constants.BuildTimestamp=$BUILD_TIMESTAMP -X github.com/hibare/Waypoint/internal/constants.CommitHash=$COMMIT_HASH" -o /bin/waypoint ./main.go

# ================== Build Frontend ================== #
FROM node:25-alpine AS frontend-build

WORKDIR /app

ARG VERSION=unknown
ARG BUILD_TIMESTAMP=unknown
ARG COMMIT_HASH=unknown

ENV VITE_VERSION=$VERSION
ENV VITE_BUILD_TIMESTAMP=$BUILD_TIMESTAMP
ENV VITE_COMMIT_HASH=$COMMIT_HASH

COPY web/package.json web/pnpm-lock.yaml* ./

# hadolint ignore=DL3016
RUN npm install -g pnpm && pnpm install

COPY web/ ./
RUN pnpm run build

# ================== Build Final Image ================== #
# hadolint ignore=DL3006
FROM alpine

# ensure consistent shell in final stage too
SHELL ["/bin/ash", "-o", "pipefail", "-c"]

ENV WAYPOINT_SERVER_LISTEN_ADDR=0.0.0.0

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=build /bin/waypoint /bin/waypoint

COPY --from=build /usr/local/bin/healthcheck /bin/healthcheck

COPY --from=frontend-build /app/dist /app/dist

HEALTHCHECK \
    --interval=30s \
    --timeout=3s \
    CMD ["healthcheck", "--url", "http://localhost:5000/ping"]

EXPOSE 5000

CMD ["/bin/waypoint", "serve"]
