FROM golang:1.21.5-alpine AS base

# ================== Build App ================== #
FROM base AS build

ARG VERSION=0.0.0

# Install healthcheck cmd
RUN apk update \
    && apk add curl \
    && apk add cosign \
    && curl -sfL https://raw.githubusercontent.com/hibare/go-docker-healthcheck/main/install.sh | sh -s -- -d -v -b /usr/local/bin

WORKDIR /src/

COPY . /src/

RUN apk --no-cache add ca-certificates

RUN CGO_ENABLED=0 go build -ldflags "-X github.com/hibare/GoGeoIP/cmd.Version=$VERSION" -o /bin/go_geo_ip ./main.go

# ================== Build Final Image ================== #
FROM alpine

ENV API_LISTEN_ADDR 0.0.0.0

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=build /bin/go_geo_ip /bin/go_geo_ip

COPY --from=build /usr/local/bin/healthcheck /bin/healthcheck

HEALTHCHECK \
    --interval=30s \
    --timeout=3s \
    CMD ["healthcheck", "--url", "http://localhost:5000/api/v1/health/"]

EXPOSE 5000

CMD ["/bin/go_geo_ip", "serve"]