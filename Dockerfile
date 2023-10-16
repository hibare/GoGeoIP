FROM golang:1.21.3-alpine AS base

# Build golang healthcheck binary
FROM base AS healthcheck

ARG VERSION=0.1.0

RUN wget -O - https://github.com/hibare/go-docker-healthcheck/archive/refs/tags/v${VERSION}.tar.gz |  tar zxf -

WORKDIR /go/go-docker-healthcheck-${VERSION}

RUN CGO_ENABLED=0 go build -o /bin/healthcheck

# Build main app
FROM base AS build

ARG GIT_VERSION_TAG=dev

WORKDIR /src/

COPY . /src/

RUN apk --no-cache add ca-certificates

RUN CGO_ENABLED=0 go build -ldflags "-X github.com/hibare/GoGeoIP/cmd.Version=$GIT_VERSION_TAG" -o /bin/go_geo_ip ./main.go

# Generate final image
FROM alpine

ENV API_LISTEN_ADDR 0.0.0.0

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=build /bin/go_geo_ip /bin/go_geo_ip

COPY --from=healthcheck /bin/healthcheck /bin/healthcheck

HEALTHCHECK \
    --interval=30s \
    --timeout=3s \
    CMD ["healthcheck","http://localhost:5000/api/v1/health/"]

EXPOSE 5000

CMD ["/bin/go_geo_ip", "serve"]