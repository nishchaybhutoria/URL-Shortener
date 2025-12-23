# Go-Redis URL Shortener

A high-performance, containerized URL shortener API built with Go (Fiber) and Redis. This project implements strict rate limiting, custom short URLs, and expiration logic, fully packaged with Docker Compose.

## Features

* **Fast & Lightweight:** Built on Fiber, an Express-inspired web framework for Go.
* **Rate Limiting:** IP-based rate limiting (10 requests/30 mins) using Redis TTL.
* **Data Persistence:** Short links are saved to disk via Redis persistence.
* **Custom Short Links:** Users can specify their own custom short IDs.
* **Link Expiry:** Links automatically expire after a set duration.
* **Dockerized:** Full microservices setup with separate API and DB containers.
* **Container Networking:** API and DB communicate via a private Docker network.

## Tech Stack

* **Language:** Golang 1.25
* **Framework:** Fiber v3
* **Database:** Redis (Alpine Image)
* **DevOps:** Docker, Docker Compose, Make
* **Libraries:** go-redis/v9, asaskevich/govalidator, google/uuid

## Project Structure

```text
.
├── api/                # Go Application Logic
│   ├── routes/         # API Route Handlers
│   ├── db/             # Database Connection Logic
│   ├── helpers/        # Validation Helpers
│   ├── main.go         # Entry Point
│   ├── Dockerfile      # API Container Config
│   └── .env.example    # Example Environment Variables
├── db/                 # Database Configuration
│   └── Dockerfile      # Redis Container Config
├── data/               # Persistent storage (mounted volume)
├── docker-compose.yml  # Container orchestration
└── Makefile            # Quick command shortcuts

```

## Getting Started

### Prerequisites

* Docker & Docker Compose
* Make (optional)

### 1. Installation

Clone the repository:

```bash
git clone [https://github.com/nishchaybhutoria/URL-Shortener.git](https://github.com/nishchaybhutoria/URL-Shortener.git)
cd url-shortener

```

### 2. Configuration

Set up your environment variables. An example file is provided in the api directory.

```bash
cp api/.env.example api/.env

```

Ensure api/.env contains the following configuration:

```ini
DB_ADDR=db:6379
DB_PASS=
APP_PORT=3000
API_QUOTA=10
DOMAIN=localhost:3000
```

Note: Docker Compose is configured to read the environment variables directly from api/.env.

### 3. Run the Application

Use the included Makefile to spin up the cluster:

```bash
make run
```

Alternatively, using Docker Compose directly:

```bash
docker-compose up -d --build
```

The server will start on http://localhost:3000.

## API Documentation

### 1. Shorten a URL

**Endpoint:** POST /api/v1

**Request Body:**

```json
{
  "url": "[https://github.com/nishchaybhutoria](https://github.com/nishchaybhutoria)",
  "short": "git",
  "expiry": 24
}
```

* **url:** The destination URL (required).
* **short:** Custom short ID (optional).
* **expiry:** Expiration time in hours (optional, default 24).

**Curl Command:**

```bash
curl -X POST -H "Content-Type: application/json" \
-d '{"url": "[https://google.com](https://google.com)", "short": "goo", "expiry": 24}' \
http://localhost:3000/api/v1
```

**Response:**

```json
{
  "url": "[https://google.com](https://google.com)",
  "short": "localhost:3000/goo",
  "expiry": 24,
  "rate_limit": "9",
  "rate_limit_reset": 29
}
```

### 2. Resolve (Redirect) a URL

**Endpoint:** GET /:url

Visit http://localhost:3000/goo in your browser.

* **301 Moved Permanently:** Redirects to the original URL.
* **404 Not Found:** If the ID does not exist or has expired.
* **503 Service Unavailable:** If the domain is blocked or rate limit exceeded.

## System Design

The application uses two Redis Databases to separate concerns:

1. **Database 0 (Storage):** Stores the mapping ShortID -> LongURL. It handles data persistence and checks for ID collisions.
2. **Database 1 (Rate Limiting):** Stores UserIP -> RequestCount. It uses Redis TTL to enforce a fixed window rate limit strategy.

## Useful Commands

| Command | Description |
| --- | --- |
| make run | Build and start the app in the background |
| make logs | View live logs from API and Redis |
| make stop | Stop all containers |
| make clean | Stop containers and remove orphans |
| make reset | Wipe all database data and restart |
