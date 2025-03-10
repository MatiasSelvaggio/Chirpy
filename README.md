# Chirpys API Documentation

## Installation

```sh
# Clone the repository
git clone https://github.com/youruser/chirpys.git
cd chirpys

# Build and run the API
go build -o chirpys
./chirpys
```

### Requirements

- Go 1.18 or higher
- PostgreSQL database configured
- Docker and Docker Compose

### Configuration

The API uses PostgreSQL and requires a `.env` file for configuration. Make sure to create a `.env` file in the root of the project with the following content:

```env
PORT=:8080
FILE_PATH_ROOT=.
DB_URL="postgres://admin:admin@localhost:5433/chirpy?sslmode=disable"
PLATFORM="dev"
JWT_SECRET_KEY="56kZfajexS2i9DMi8jrR4V2Lf4icVUAvWczLAPm+SmfMQEYz25gzDaKRyOT9hoYsPjzdnIdVdRtJk4v3eqBjsg=="
POLKA_KEY="f271c81ff7084ee5b99a5091b42d486e"
```

### Start the Database with Docker

If you do not have PostgreSQL installed, you can start it with Docker Compose:

```sh
docker-compose up -d
```

## API Usage

The following endpoints are available in the API.

### GET /api/healthz

This endpoint checks the health status of the API.

#### Request

```http
GET /api/healthz HTTP/1.1
Host: localhost:8080
```

#### Response

```json
{
  "status": "ok"
}
```

#### Status Codes

- `200 OK` - The API is working correctly.
- `500 Internal Server Error` - An internal error occurred in the API.

### POST /api/users

This endpoint creates a new user.

#### Request

```http
POST /api/users HTTP/1.1
Host: localhost:8080
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword"
}
```

#### Response

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "created_at": "2025-03-10T12:00:00Z",
  "updated_at": "2025-03-10T12:00:00Z",
  "email": "user@example.com",
  "is_chirpy_red": false
}
```

#### Status Codes

- `201 Created` - The user was successfully created.
- `400 Bad Request` - The email or password is missing.
- `409 Conflict` - The email is already in use.
- `500 Internal Server Error` - A system error occurred.

### PUT /api/users

This endpoint allows updating the user's email and password. A JWT token must be provided in the `Authorization` header with the `Bearer` prefix.

#### Request

```http
PUT /api/users HTTP/1.1
Host: localhost:8080
Content-Type: application/json
Authorization: Bearer <your-jwt-token>

{
  "email": "newuser@example.com",
  "password": "newsecurepassword"
}
```

#### Response

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "created_at": "2025-03-10T12:00:00Z",
  "updated_at": "2025-03-10T12:30:00Z",
  "email": "newuser@example.com",
  "is_chirpy_red": false
}
```

#### Status Codes

- `200 OK` - The user was successfully updated.
- `400 Bad Request` - The email or password is missing.
- `401 Unauthorized` - No valid JWT token was provided.
- `409 Conflict` - The new email is already in use.
- `500 Internal Server Error` - A system error occurred.

### POST /api/login

This endpoint is used to obtain a JWT token for authenticating API requests.

#### Request

```http
POST /api/login HTTP/1.1
Host: localhost:8080
Content-Type: application/json

{
  "email": "walt@breakingbad.com",
  "password": "securepassword"
}
```

#### Response

```json
{
  "id": "9db97406-65f4-4e44-92d0-3998abe4c09e",
  "created_at": "2025-03-10T21:37:31.98102Z",
  "updated_at": "2025-03-10T21:37:31.98102Z",
  "email": "walt@breakingbad.com",
  "is_chirpy_red": false,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaGlycHkiLCJzdWIiOiI5ZGI5NzQwNi02NWY0LTRlNDQtOTJkMC0zOTk4YWJlNGMwOWUiLCJleHAiOjE3NDE2NDc4MDAsImlhdCI6MTc0MTY0NDIwMH0.ktNIuvEpp1ygyb1r92Qzqj-QM0_8jnP87R3VBvgZL2U",
  "refresh_token": "02770f807ce3a1c06a686df93d3e377e903a41fffa7a0c220e26e639f9ce6972"
}
```

#### Token Information

- The JWT token is valid for **1 hour**.
- The refresh token allows generating a new JWT token for **60 days** or until revoked.

#### Status Codes

- `200 OK` - Login successful, JWT and refresh token returned.
- `400 Bad Request` - The email or password is missing.
- `401 Unauthorized` - Invalid email or password.
- `500 Internal Server Error` - A system error occurred.

### POST /api/refresh

This endpoint is used to refresh the JWT token using a refresh token.

#### Request

```http
POST /api/refresh HTTP/1.1
Host: localhost:8080
Authorization: Bearer <your-refresh-token>
```

#### Response

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaGlycHkiLCJzdWIiOiI5ZGI5NzQwNi02NWY0LTRlNDQtOTJkMC0zOTk4YWJlNGMwOWUiLCJleHAiOjE3NDE2NDgwODIsImlhdCI6MTc0MTY0NDQ4Mn0.v3zAE0k7yZBPetwQsJgoWo-qCv658FQTlPlersSAais"
}
```

#### Status Codes

- `200 OK` - Token refreshed successfully.
- `401 Unauthorized` - Invalid or expired refresh token.
- `500 Internal Server Error` - A system error occurred.

