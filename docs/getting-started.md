# Getting Started with Waypoint

Waypoint is a self-hosted IP geolocation service. This guide will help you get up and running.

## Prerequisites

- **Docker & Docker Compose**: For running the service
- **MaxMind License Key**: Get a free account at [maxmind.com](https://www.maxmind.com)

## Installation & Setup

### Docker-based Setup (Recommended)

1. **Clone the repository**:

   ```bash
   git clone https://github.com/hibare/Waypoint.git
   cd Waypoint
   ```

2. **Initialize environment**:

   ```bash
   cp .env.example .env
   ```

3. **Configure environment**:
   Edit `.env` and set:
   - `POSTGRES_PASSWORD` - A secure password for PostgreSQL
   - `WAYPOINT_MAXMIND_LICENSE_KEY` - Your MaxMind license key

4. **Start Waypoint**:

   ```bash
   docker compose up -d
   ```

### Manual Setup (For local development)

1. **Clone the repository**:

   ```bash
   ://github.com/hibare/Way   git clone httpspoint.git
   cd Waypoint
   ```

2. **Backend Dependencies**:

   ```bash
   go mod download
   ```

3. **Frontend Dependencies**:
   ```bash
   cd web && pnpm install
   ```

## Configuration

Waypoint uses a YAML configuration file. An example is provided in `config.example.yaml`.

1. Copy the example config:

   ```bash
   cp config.example.yaml config.yaml
   ```

2. Edit `config.yaml` with your settings:
   - Database credentials (optional, uses in-memory by default)
   - MaxMind license key (required)
   - OIDC settings (optional)

## Running Waypoint

### Docker

```bash
docker compose up -d
```

The service will be available at:

- Web UI: http://localhost:5000
- API: http://localhost:5000/api/v1

### CLI

```bash
# Start server
waypoint serve

# Download MaxMind databases
waypoint maxmind download

# Lookup IP
waypoint lookup 8.8.8.8

# Run database migrations
waypoint db migrate
```

---

[← Back to Overview](../README.md)
