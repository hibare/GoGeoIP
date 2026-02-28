# API Reference

Waypoint provides a REST API for IP geolocation lookups.

## Base URL

```txt
http://localhost:5000/api/v1
```

## Authentication

### API Key Authentication

Include your API key in the `Authorization` header:

```bash
curl -H "Authorization: YOUR_API_KEY" http://localhost:5000/api/v1/ip/8.8.8.8
```

### Cookie Authentication

When OIDC is enabled, you can authenticate via browser cookies after logging in through the web UI.

## Endpoints

### Health Check

Check if the service is running.

**Endpoint:** `GET /ping`

**Response:**

```json
{
  "ok": true
}
```

### Get My IP

Get the public IP address of the requesting client.

**Endpoint:** `GET /api/v1/ip`

**Response:**

```json
{
  "city": "Mountain View",
  "country": "United States",
  "continent": "North America",
  "iso_country_code": "US",
  "iso_continent_code": "NA",
  "is_anonymous_proxy": false,
  "is_satellite_provider": false,
  "timezone": "America/Los_Angeles",
  "latitude": 37.4223,
  "longitude": -122.0848,
  "asn": 15169,
  "organization": "GOOGLE",
  "ip": "8.8.8.8"
}
```

### Lookup IP Address

Get Geo location information for a specific IP address.

**Endpoint:** `GET /api/v1/ip/{ip}`

**Parameters:**

- `ip` - IP address to lookup (IPv4 or IPv6)

**Headers:**

- `Authorization` - API key (required for protected endpoints)

**Response:**

```json
{
  "city": "Mountain View",
  "country": "United States",
  "continent": "North America",
  "iso_country_code": "US",
  "iso_continent_code": "NA",
  "is_anonymous_proxy": false,
  "is_satellite_provider": false,
  "timezone": "America/Los_Angeles",
  "latitude": 37.4223,
  "longitude": -122.0848,
  "asn": 15169,
  "organization": "GOOGLE",
  "ip": "8.8.8.8"
}
```

### List API Keys

List all API keys for the authenticated user.

**Endpoint:** `GET /api/v1/api-keys`

**Headers:**

- `Authorization` - Cookie or API key

**Response:**

```json
[
  {
    "id": "uuid",
    "name": "My API Key",
    "state": "active",
    "scopes": [],
    "expires_at": null,
    "created_at": "2024-01-01T00:00:00Z",
    "last_used_at": "2024-01-02T00:00:00Z"
  }
]
```

### Create API Key

Create a new API key.

**Endpoint:** `POST /api/v1/api-keys`

**Headers:**

- `Authorization` - Cookie or API key

**Request Body:**

```json
{
  "name": "My API Key",
  "expires_at": "2024-12-31T23:59:59Z"
}
```

**Response:** Returns the newly created API key (shown only once).

### Revoke API Key

Revoke an API key.

**Endpoint:** `POST /api/v1/api-key/{id}/revoke`

**Parameters:**

- `id` - API key UUID

**Headers:**

- `Authorization` - Cookie or API key

### Delete API Key

Delete an API key permanently.

**Endpoint:** `DELETE /api/v1/api-key/{id}`

**Parameters:**

- `id` - API key UUID

**Headers:**

- `Authorization` - Cookie or API key

### Get Current User

Get information about the authenticated user.

**Endpoint:** `GET /api/v1/auth/me`

**Headers:**

- `Authorization` - Cookie or API key

**Response:**

```json
{
  "id": "user-uuid",
  "email": "user@example.com",
  "name": "John Doe",
  "picture": "https://example.com/avatar.jpg"
}
```

## Rate Limits

There are no rate limits on the API.

## Error Responses

### 401 Unauthorized

```json
{
  "error": "unauthorized"
}
```

### 404 Not Found

```json
{
  "error": "not found"
}
```

### 500 Internal Server Error

```json
{
  "error": "internal server error"
}
```
