# GoGeoIP

[![Go Report Card](https://goreportcard.com/badge/github.com/hibare/GoGeoIP)](https://goreportcard.com/report/github.com/hibare/GoGeoIP)
[![Docker Hub](https://img.shields.io/docker/pulls/hibare/go-geo-ip)](https://hub.docker.com/r/hibare/go-geo-ip)
[![Docker image size](https://img.shields.io/docker/image-size/hibare/go-geo-ip/latest)](https://hub.docker.com/r/hibare/go-geo-ip)
[![GitHub issues](https://img.shields.io/github/issues/hibare/GoGeoIP)](https://github.com/hibare/GoGeoIP/issues)
[![GitHub pull requests](https://img.shields.io/github/issues-pr/hibare/GoGeoIP)](https://github.com/hibare/GoGeoIP/pulls)
[![GitHub](https://img.shields.io/github/license/hibare/GoGeoIP)](https://github.com/hibare/GoGeoIP/blob/main/LICENSE)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/hibare/GoGeoIP)](https://github.com/hibare/GoGeoIP/releases)

A fully self-hosted Rest IP to fetch Geo information for an IP.
Geo information such as City, Country, ASN, Organization etc., details.

Rest API depends on [MaxMind](https://www.maxmind.com/en/home) City, Country & ASN Lite database.

MaxMind offers City, Country & ASN Lite databases for free for individual uses.

The program downloads City, Country and ASN databases from MaxMind and runs lookups on local DB.

## Features

- Rest API to translate IP to Geo info
- Automatically download, verify and load MaxMind lite databases on boot.
- Regularly run jobs to update MaxMind local databases.
- No API lookup limits.

## Getting Started

GoGeoIP is packaged as docker container. Docker image is available on [Docker Hub](https://hub.docker.com/r/hibare/go-geo-ip).

To use GoGeoIP, you require MaxMind license key to download lite databases. Head over to [MaxMind](https://www.maxmind.com/en/geolite2/signup?utm_source=kb&utm_medium=kb-link&utm_campaign=kb-create-account) and sign-up for a free account.

### Get MaxMind License Key

- Login to [MaxMind](https://www.maxmind.com/en/account/login)
- Click on `Manage License Keys` in left side menu.
- Click on `Generate new license key`, fill out description and click confirm.
- Save license key safely

### Docker run

```shell
docker run -it -p 5000:5000 -e DB_LICENSE_KEY=<LICENSE_KEY> -e API_LISTEN_ADDR=0.0.0.0 -e API_KEYS=test-key  hibare/go-geo-ip
```

Replace `<LICENSE_KEY>` with the license key from MaxMind.
Replace `test-key` with randomly generated API key. This is used to authenticate all IP lookup rest calls.

```
INFO[0000] Loaded config
INFO[0000] Downloading all DB files
INFO[0000] Downloading DB file, path=/tmp/GeoLite2-Country.tar.gz
INFO[0000] Scheduling DB update job
INFO[0003] Downloaded DB file, path=/tmp/GeoLite2-Country.tar.gz
INFO[0003] Downloading sha256 file, path=/tmp/GeoLite2-Country.tar.gz.sha256
INFO[0005] Downloaded sha256 file, path=/tmp/GeoLite2-Country.tar.gz.sha256
INFO[0005] Checksum validated for archive /tmp/GeoLite2-Country.tar.gz
INFO[0005] Extracting file GeoLite2-Country.mmdb from archive /tmp/GeoLite2-Country.tar.gz
INFO[0005] Extracted file GeoLite2-Country.mmdb from archive /tmp/GeoLite2-Country.tar.gz at /tmp/GeoLite2-Country.mmdb
INFO[0005] Loading new DB file data/GeoLite2-Country.mmdb
INFO[0005] New DB file loaded data/GeoLite2-Country.mmdb
INFO[0005] Downloading DB file, path=/tmp/GeoLite2-City.tar.gz
INFO[0013] Downloaded DB file, path=/tmp/GeoLite2-City.tar.gz
INFO[0013] Downloading sha256 file, path=/tmp/GeoLite2-City.tar.gz.sha256
INFO[0014] Downloaded sha256 file, path=/tmp/GeoLite2-City.tar.gz.sha256
INFO[0014] Checksum validated for archive /tmp/GeoLite2-City.tar.gz
INFO[0014] Extracting file GeoLite2-City.mmdb from archive /tmp/GeoLite2-City.tar.gz
INFO[0015] Extracted file GeoLite2-City.mmdb from archive /tmp/GeoLite2-City.tar.gz at /tmp/GeoLite2-City.mmdb
INFO[0015] Loading new DB file data/GeoLite2-City.mmdb
INFO[0015] New DB file loaded data/GeoLite2-City.mmdb
INFO[0015] Downloading DB file, path=/tmp/GeoLite2-ASN.tar.gz
INFO[0017] Downloaded DB file, path=/tmp/GeoLite2-ASN.tar.gz
INFO[0017] Downloading sha256 file, path=/tmp/GeoLite2-ASN.tar.gz.sha256
INFO[0018] Downloaded sha256 file, path=/tmp/GeoLite2-ASN.tar.gz.sha256
INFO[0018] Checksum validated for archive /tmp/GeoLite2-ASN.tar.gz
INFO[0018] Extracting file GeoLite2-ASN.mmdb from archive /tmp/GeoLite2-ASN.tar.gz
INFO[0018] Extracted file GeoLite2-ASN.mmdb from archive /tmp/GeoLite2-ASN.tar.gz at /tmp/GeoLite2-ASN.mmdb
INFO[0018] Loading new DB file data/GeoLite2-ASN.mmdb
INFO[0018] New DB file loaded data/GeoLite2-ASN.mmdb
INFO[0018] Listening for address 0.0.0.0 on port 5000
```

### Docker Compose

Create a `.env` file and copy the content from `.env.example`. Alternatively, rename `.env.example` to `.env`. Replace values of all variables in `.env` file with appropriate values.

A minimalistic docker-compose.yml file is provided in the repo. Download docker-compose.yml file.

```shell
curl https://raw.githubusercontent.com/hibare/GoGeoIP/main/docker-compose.yml -o docker-compose.yml
```

Run docker-compose.yml file

```shell
docker compose up
```

## Endpoints

1. Check health
   `[GET]` `/api/v1/health`

```shell
❯ curl http://127.0.0.1:5000/api/v1/health
```

```json
{ "ok": true }
```

2. IP Geo
   `[GET]` `/api/v1/ip/{lookup_ip}`

```shell
❯ curl -H "Authorization: test-key" http://127.0.0.1:5000/api/v1/ip/106.213.87.64
```

```json
{
  "city": {
    "city": "Los Angeles",
    "country": "United States",
    "continent": "North America",
    "iso_country_code": "US",
    "iso_continent_code": "NA",
    "is_anonymous_proxy": false,
    "is_satellite_provider": false,
    "timezone": "America/Los_Angeles",
    "latitude": 34.0544,
    "longitude": -118.2441
  },
  "asn": { "asn": 15169, "oraganization": "GOOGLE" }
}
```

## Cli

GoGeoIP also has cli commands for quick actions. Binary is `go_geo_ip`.

For docker containers prefix all commands with `docker exec -it {container_name}`

```shell
go_geo_ip -h
```

```
API to fetch Geo information for an IP

Usage:
  go_geo_ip [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  db          IP DB related commands
  geoip       Lookup Geo information for an IP
  help        Help about any command
  keys        Manage API Keys
  serve       Start API Server

Flags:
  -h, --help      help for go_geo_ip
  -v, --version   version for go_geo_ip

Use "go_geo_ip [command] --help" for more information about a command.
```

## Start API Server

```shell
go_geo_ip serve
```

## Download DB

```shell
go_geo_ip db download
```

## List API Keys

```shell
go_geo_ip keys list
```

## Version

```shell
go_geo_ip --version
```

## Environment Variables

| Variable               | Description                                          | Required | Default Value                 | Value Type                      |
| ---------------------- | ---------------------------------------------------- | -------- | ----------------------------- | ------------------------------- |
| API_LISTEN_ADDR        | IP address to bind API server                        | No       | 0.0.0.0                       | string                          |
| API_LISTEN_PORT        | Port to listen                                       | No       | 5000                          | int                             |
| API_KEYS               | Comma separated API keys to authenticated REST calls | No       | Auto generated during runtime | comma separated string          |
| DB_LICENSE_KEY         | MaxMind License key                                  | Yes      | -                             | string                          |
| DB_AUTOUPDATE          | Flag to enable/disable DB auto-update                | No       | true                          | boolean                         |
| DB_AUTOUPDATE_INTERVAL | Auto update interval.                                | No       | 24 Hours                      | Time duration (ex: 24h, 1h, 6h) |
