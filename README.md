<div align="center">
  <img src="./web/public/logo.png" alt="Waypoint Logo" width="200" height="200">

# Waypoint

_A self-hosted IP geolocation service_

[![Go Report Card](https://goreportcard.com/badge/github.com/hibare/Waypoint)](https://goreportcard.com/report/github.com/hibare/Waypoint)
[![GitHub issues](https://img.shields.io/github/issues/hibare/Waypoint)](https://github.com/hibare/Waypoint/issues)
[![GitHub pull requests](https://img.shields.io/github/issues-pr/hibare/Waypoint)](https://github.com/hibare/Waypoint/pulls)

</div>

Waypoint is a self-hosted IP geolocation service that provides geographic information for any IP address including City, Country, ASN, Organization, and more using MaxMind GeoIP databases.

## 🚀 Key Features

- **REST API**: Fast and reliable IP geolocation lookup
- **Web UI**: Clean and intuitive dashboard for IP lookups and history
- **API Key Management**: Secure API key generation and revocation
- **User Authentication**: OIDC support for secure access
- **Automatic Updates**: MaxMind database auto-update
- **API First**: Fully featured REST API for seamless integration

## 📁 Project Structure

```text
Waypoint/
├── cmd/             # Executable entry points (Server, CLI)
├── internal/        # Core business logic
│   ├── config/      # Configuration management
│   ├── db/          # Database layer
│   ├── handlers/    # API request handlers
│   └── maxmind/     # MaxMind client
├── web/             # Vue.js frontend application
└── compose.yml      # Docker Compose configuration
```

## 📖 Documentation

- [**Getting Started**](./docs/getting-started.md): Installation, configuration, and first steps.
- [**API Reference**](./docs/api.md): Overview of available REST endpoints.
- [**Configuration**](./config.example.yaml): Configuration options and environment variables.

## License

[MIT](LICENSE)
