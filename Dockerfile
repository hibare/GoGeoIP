ARG GOLANG_VERSION=1.24.0

FROM golang:${GOLANG_VERSION}-alpine AS base

# ================== Build App ================== #
FROM base AS build

# prefer ash on Alpine so -o pipefail behaves as expected
SHELL ["/bin/ash", "-o", "pipefail", "-c"]

ARG VERSION=0.0.0

# Install healthcheck cmd
# hadolint ignore=DL3018
RUN apk update \
    && apk add --no-cache curl cosign  \
    && rm -rf /var/cache/apk/* \
    && curl -sfL https://raw.githubusercontent.com/hibare/go-docker-healthcheck/main/install.sh | sh -s -- -d -v -b /usr/local/bin

WORKDIR /src/

COPY . /src/

# hadolint ignore=DL3018
RUN apk --no-cache add ca-certificates \
    && CGO_ENABLED=0 go build -ldflags "-X github.com/hibare/GoGeoIP/cmd.Version=$VERSION" -o /bin/go_geo_ip ./main.go

# ================== Build Final Image ================== #
# hadolint ignore=DL3006
FROM alpine

# ensure consistent shell in final stage too
SHELL ["/bin/ash", "-o", "pipefail", "-c"]

ENV SERVER_LISTEN_ADDR=0.0.0.0

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=build /bin/go_geo_ip /bin/go_geo_ip

COPY --from=build /usr/local/bin/healthcheck /bin/healthcheck

HEALTHCHECK \
    --interval=30s \
    --timeout=3s \
    CMD ["healthcheck", "--url", "http://localhost:5000/api/v1/health/"]

EXPOSE 5000

CMD ["/bin/go_geo_ip", "serve"]
