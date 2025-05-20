<!-- logo -->
<p align="center">
  <img src="email/assets/logo.png" alt="URL Shortener Logo"/>
</p>

# URL Shortener API

A production-grade, extensible backend for URL shortening.  
Designed for reliability, observability, and maintainability with Go, PostgreSQL, Redis, Prometheus, Grafana, Alertmanager, and Traefik.  
No frontend included – this is a robust, scalable **service layer**.

---

## Key Features

- **Shorten URLs** with optional custom codes
- **User registration, login, JWT authentication**
- **Email verification** on signup, password reset with branded HTML email (with static images)
- **Admin tools**: full user/URL management, only for admins
- **Click batching:** Clicks stored in PostgreSQL until a threshold, then hot URLs (those with significant hits) are promoted to Redis for fast access; periodic "flusher job" moves counts back to DB and clears Redis.
- **Rate limiting:** Per-IP, via Traefik, fully configurable
- **Health job:** Scheduled background checks on database, Redis, mail; serves `/health` endpoint as JSON and feeds status to Prometheus (for alerting)
- **Full monitoring/alerting:** Prometheus metrics, custom alert rules, dashboards (Grafana), email alerts (Alertmanager)
- **Observalibity:** Each system (app, db, flusher, health, mail, server) logs to a separate file under `/logs`
- **Strict security:** JWT Bearer auth, admin-only endpoints, IP whitelisting & Basic Auth for admin/metrics/dashboards
- **Makefile** for local/devops convenience
- **CI/CD out of the box** (GitHub Actions + Docker Hub auto-publish)
- **Traefik-first:** Designed for containerized, HTTPS, and secure reverse proxy use

---

## API Documentation

- **OpenAPI/Swagger:**
    - Auto-generated with [swaggo/swag](https://github.com/swaggo/swag)
    - Available at [API Doc](https://cemakan.com.tr/api/docs/index.html)
    - All endpoints, models, and authentication methods are fully documented

---

## Observability & Health System

- **/health endpoint**: Returns live JSON with DB/Redis/Mail status (used by both users and Prometheus scraping)
- **Health Job:** Background process regularly checks all core services (DB, Redis, Mail). If any fail, `/health` and Prometheus `*_up` metrics reflect this.
- **Flusher job:** Periodically flushes click stats from Redis to PostgreSQL, deleting from Redis after successful DB write.
    - Only URLs with real/hot traffic are promoted to Redis. Cold/rarely-hit links are fetched straight from PostgreSQL.
- **Alerting (Prometheus + Alertmanager):**
    - **System Alerts:** CPU > 80%/RAM > 85%, Redis/DB/Mail down
    - **App Alerts:** Slow redirects, low/high request rate, flusher stalls
    - **Alertmanager:** Sends notifications to configured email (`TEAM_EMAIL`)
    - All alert rules and intervals are configurable via Prometheus rules files.
- **Grafana Dashboards:**
    - Only available to whitelisted IPs with Basic Auth
    - Access at `/grafana` (see docker-compose and Traefik setup)

---

## Security

- **JWT Bearer token** for all protected endpoints
- **IP whitelisting & Basic Auth** for:
    - `/metrics`, `/api/docs`, `/grafana`, `/prometheus`, `/alert`
- **Password hashing** and secure admin/user separation
- **Admin seeding:** If no admin exists, a first one is auto-created from `.env` at boot

---

## Automation & DevOps

- **CI/CD:**
    - Build, test, and publish Docker images to [Docker Hub](https://hub.docker.com/r/cemakan/url-shortener) via GitHub Actions
    - See `.github/workflows/ci-cd.yml`
- **Makefile:**
    - Fast local Docker, migration, and devops commands
- **Traefik**: All routing, SSL, path-based access, rate limits via Traefik + ENV config
- **Mail alerts:** Alertmanager triggers emails for critical events to your team

---

## Deployment & Requirements

- **Docker Compose**:
    - Includes everything: app, PostgreSQL, Redis, Prometheus, Grafana, Alertmanager, node/postgres exporters
- **Requires:**
    - External **SMTP server** for all email features (register, reset)
    - **Traefik** as a reverse proxy (HTTPS, rate limiting, path-prefix)
- **Example .env:**
    - See `.env.example` for all required configuration

---

## Project Structure

- `cmd/`        — Application entrypoint
- `config/`     — Environment and JWT config
- `docs/`       — Swagger docs (auto-generated)
- `email/`      — Email templates & static assets
- `internal/`   — All core logic, handlers, middleware, domain, health, infrastructure, services
- `logs/`       — Separate log files per subsystem
- `prometheus/` — Prometheus and alert rules config
- `traefik/`    — Example reverse proxy configs
- `Makefile`, `docker-compose.yml`, etc.

---

## API Quick Reference

- `POST   /api/register`         – Register (sends email confirmation)
- `POST   /api/login`            – Obtain JWT token
- `GET    /api/health`           – Health info (DB, Redis, Mail)
- `POST   /api/shorten`          – Shorten a URL (auth required)
- `GET    /api/me`               – User info (auth required)
- `GET    /api/my/urls`          – All user URLs (auth required)
- `GET    /api/{code}`           – Redirect by code (public)
- `GET    /api/docs/index.html`  – Full API reference
- `GET    /api/admin/users`      – List users (admin only)
- `DELETE /api/admin/users/{id}` – Remove user (admin only)

**Admin only:**
- `GET    /admin/users`
- `DELETE /admin/users/{id}`

---

## Monitoring & Alerting

- **System Metrics:**
    - CPU, memory, disk, container health
    - Redis/DB/Mail up/down status
- **Application Metrics:**
    - Click rates, 95th percentile redirect latency, request rates
    - Flusher job health (clicks flushed, failure count, stalls)
- **Alerts:**
    - Configurable in `prometheus/rules.yml`
    - Mail notifications via Alertmanager (configurable in `.env`)
- **Grafana Dashboards:**
    - Visualize all metrics (protected by Basic Auth + IP whitelist)

---

## Usage

### 1. Clone & configure

```sh
git clone https://github.com/CemAkan/url-shortener.git
cd url-shortener
cp .env.example .env
# Edit .env with your DB, Redis, SMTP, JWT, admin, etc.
```

# TODO

- [ ] QR code support for short links
- [ ] User-selectable cache/flush policy (.env switch)
- [ ] “Cold” URL Redis eviction tuning
- [ ] Self-service API keys & integrations
- [ ] More admin tools & exports
- [ ] Optional “delayed flush” for burst traffic
  