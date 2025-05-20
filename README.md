<!-- logo -->
<p align="center">
  <img src="email/assets/logo.png" alt="URL Shortener Logo" />
</p>

<h1 align="center">URL Shortener API</h1>
<p align="center">
  A production-grade, secure, and extensible backend for high-performance URL shortening.
</p>

---

## Overview

URL Shortener is a containerized, enterprise-ready backend API for managing short URLs, built with Go and following clean architecture principles. It integrates PostgreSQL, Redis, Prometheus, Grafana, Alertmanager, and Traefik. The system emphasizes scalability, observability, security, and operational maintainability.

This project includes no frontend. It is designed to be integrated into existing platforms (web, mobile, CLI) or to run as a dedicated microservice.

---

## Key Features

- REST API for shortening URLs (with optional custom aliases)
- Secure user registration, login, and JWT-based authentication
- Email verification and password reset using branded HTML templates
- Admin-only endpoints for managing users and URLs
- Click batching system using PostgreSQL and Redis
- Real-time health check API and periodic background health verifier
- Per-IP rate limiting via Traefik middleware
- Prometheus metrics and Grafana dashboards
- Email alerts for system and application-level anomalies
- Per-component structured logging
- GitHub Actions-based CI/CD and Docker Hub deployment
- Secure HTTPS routing using Traefik with automatic TLS

---

## Architecture

| Component           | Description                                             |
|--------------------|---------------------------------------------------------|
| `url-shortener`     | Go-based REST API                                       |
| `postgres`          | Persistent relational database                          |
| `redis`             | Caching layer for hot URLs                              |
| `prometheus`        | Metrics collection and alert rule engine                |
| `grafana`           | Dashboard visualization                                 |
| `alertmanager`      | Alert delivery engine (email notifications)             |
| `node_exporter`     | Host-level metrics (CPU, RAM, Disk)                     |
| `postgres_exporter` | PostgreSQL performance metrics for Prometheus           |
| `traefik`           | HTTPS reverse proxy, rate limiter, and router           |

All services are managed via Docker Compose and isolated with internal Docker networks.

---

## API Documentation

The API is fully documented using Swagger (OpenAPI 3.0) via [swaggo/swag](https://github.com/swaggo/swag).

**Live Docs:**  
[API Swagger Live Documentation](https://cemakan.com.tr/api/docs/index.html)

The documentation includes:
- All available endpoints and methods
- Field schemas and validation rules
- Authentication structure (JWT)
- Example requests and responses
- Error responses and status codes

---

## Security

- JWT Bearer tokens for all authenticated endpoints
- IP allowlisting and Basic Auth for:
  - `/metrics`
  - `/grafana`
  - `/prometheus`
  - `/alert`
  - `/api/docs`
- Passwords stored using bcrypt hashing
- Auto-seeding of an admin account on first launch if none exists
- All traffic securely routed via Traefik with TLS (Let's Encrypt or manual certs)

---

## Monitoring & Health

- `/health` endpoint returns live JSON status of:
  - PostgreSQL
  - Redis
  - SMTP
- Background health checker runs periodically
- Prometheus scrapes metrics from:
  - Application
  - Redis
  - PostgreSQL
  - Node exporter
- Alertmanager triggers alerts via email based on thresholds or failures
- Grafana dashboards provide visual monitoring for all systems

### Default Alerts

- CPU usage above 80%
- Memory usage above 85%
- Redis, PostgreSQL, or SMTP service down
- API redirect latency degradation
- Request throughput outside configured bounds
- Flusher failure or unexpected idle time

---

## API Endpoints

| Method | Endpoint                    | Access        | Description                              |
|--------|-----------------------------|---------------|------------------------------------------|
| POST   | `/api/register`             | Public        | Register a new user                      |
| POST   | `/api/login`                | Public        | Obtain JWT access token                  |
| GET    | `/api/me`                   | Authenticated | Retrieve logged-in user profile          |
| POST   | `/api/shorten`              | Authenticated | Create a new shortened URL               |
| GET    | `/api/{code}`               | Public        | Resolve and redirect from short code     |
| GET    | `/api/my/urls`              | Authenticated | List URLs created by the current user    |
| GET    | `/api/admin/users`          | Admin Only    | List all registered users                |
| DELETE | `/api/admin/users/{id}`     | Admin Only    | Remove user by ID                        |
| GET    | `/api/docs/index.html`      | Protected     | View full API documentation              |
| GET    | `/health`                   | Public        | Get real-time system health report       |

---

## Reverse Proxy & Routing

All routes are proxied securely via Traefik with rate limiting and middleware protections.

| Path Prefix     | Destination        | Access Protection            |
|-----------------|--------------------|------------------------------|
| `/api`          | url-shortener      | JWT-based, rate limited      |
| `/grafana`      | grafana            | IP allowlist + Basic Auth    |
| `/prometheus`   | prometheus         | IP allowlist + Basic Auth    |
| `/alert`        | alertmanager       | IP allowlist + Basic Auth    |
| `/metrics`      | exporters + app    | IP allowlist + Basic Auth    |
| `/health`       | url-shortener      | Public                       |

Rate limit settings are configured in `.env` and applied via Traefik middlewares.

---

## File & Folder Structure

| Path                  | Description                                    |
|-----------------------|------------------------------------------------|
| `cmd/`                | Application entrypoint                         |
| `config/`             | Configuration loading (env, JWT, SMTP, etc.)  |
| `internal/`           | Business logic, services, handlers             |
| `email/`              | Static HTML templates and assets               |
| `docs/`               | Auto-generated Swagger specs                   |
| `prometheus/`         | Alerting rules and scrape configs              |
| `traefik/`            | Traefik dynamic config (middlewares, routers)  |
| `logs/`               | Per-component log output                       |
| `Makefile`            | DevOps helper commands                         |
| `.env.example`        | Sample environment configuration               |
| `docker-compose.yml`  | Orchestration for all services                 |

---

## Deployment Requirements

- Docker Engine 20.10+
- Docker Compose v2.20+
- Public domain name (e.g. `cemakan.com.tr`)
- SMTP credentials (username/password)
- Public access to port 443 for Let’s Encrypt TLS

---

## Deployment Steps

```bash
git clone https://github.com/CemAkan/url-shortener.git
cd url-shortener

cp .env.example .env
nano .env  # configure DB, Redis, JWT_SECRET, SMTP, etc.

docker compose up -d --build
```


---

## Accessible Routes After Deployment

| Service            | URL                                     |
|--------------------|------------------------------------------|
| API                | `https://yourdomain.com/api`             |
| Swagger Docs       | `https://yourdomain.com/api/docs`        |
| Metrics            | `https://yourdomain.com/metrics`         |
| Grafana Dashboard  | `https://yourdomain.com/grafana`         |
| Prometheus UI      | `https://yourdomain.com/prometheus`      |
| Alertmanager       | `https://yourdomain.com/alert`           |
| Healthcheck        | `https://yourdomain.com/health`          |

---

## CI/CD & Automation

- **GitHub Actions**:
  - On every push, Docker images are automatically built, tested, and published.
- **Docker Image**:
  - Available at [`cemakan/url-shortener`](https://hub.docker.com/r/cemakan/url-shortener)
- **Makefile Features**:
  - Local testing utilities
  - Database migration automation
  - Service cleanup routines
  - Centralized log rotation/archive commands
- **Log Structure**:
  - Logs are stored in `/logs/{component}.log`, organized per subsystem.

---

## 📝 Project Goals

| Feature                                  | Status |
|------------------------------------------|--------|
| Full API with JWT and Admin Role         | ✅     |
| SMTP Email Verification                  | ✅     |
| Redis Click Cache                        | ✅     |
| Prometheus Monitoring                    | ✅     |
| Grafana Dashboards                       | ✅     |
| Alertmanager Email Alerts                | ✅     |
| Rate Limiting via Traefik Middleware     | ✅     |
| IP Whitelist & Basic Auth                | ✅     |
| Healthcheck Endpoint & Verifier Job      | ✅     |
| CI/CD with GitHub Actions & Docker Hub   | ✅     |
| Per-component Log Architecture           | ✅     |
| QR Code Support for Short URLs           | ⏳     |
| Delayed Cache Flushing Strategy          | ⏳     |
| Self-service API Key Management          | ⏳     |
| CSV Export for Admin Panels              | ⏳     |
| Multi-tenant Domain Routing Support      | ⏳     |
| Role-based Access Control (RBAC)         | ⏳     |
| Cold Key Eviction (LRU/TTL in Redis)     | ⏳     |
| Real-time Click Analytics Visualization  | ⏳     |

> ✅ Implemented  ⏳ Planned/In Progress
---

## License

This project is licensed under the **MIT License**.  
See the [LICENSE](./LICENSE) file for complete details.

---

## Maintainer

**Cem Akan**  
Docker Hub: [cemakan/url-shortener](https://hub.docker.com/r/cemakan/url-shortener)

For contributions, feature requests, or bug reports, please [open an issue](https://github.com/CemAkan/url-shortener/issues) or submit a pull request via GitHub.

---