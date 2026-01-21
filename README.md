# Spinner V2 - HUT DP TASPEN 36

## Overview
This repository contains the full web experience for the HUT DP TASPEN 36 event, including invitations, portal lookup by NIP, QR scan attendance, realtime monitoring, doorprize spinner, admin panel, and winner management. The app is a Go server that renders HTML templates, serves static assets, and exposes JSON APIs backed by MySQL.

## Key Features
- Invitation page with QR code generation for attendance confirmation.
- Portal page for redirecting employees to their invitation link based on NIP.
- QR scanner page with per-scanner identifier (1-3) and attendance recording.
- Realtime monitoring page that shows the last three scans per scanner via WebSocket.
- Doorprize spinner with persisted winners and rules enforcement.
- Admin panel for attendance, guests, exports, and winner management.

## Tech Stack
- Go (chi router, sqlx, MySQL)
- Goose for database migrations
- WebSocket for realtime monitor updates
- HTML templates + static assets under `internal/templates/assets`

## Routes
### Public Pages
- `/` Invitation page (uses `?name=...`)
- `/portal` NIP portal redirect page
- `/scan` QR scanner page (requires `scanner_id` input 1-3)
- `/monitor` Realtime monitoring page

### Protected Pages (Basic Auth)
- `/spinner` Doorprize spinner
- `/admin` Admin dashboard

### Static Assets
- `/public/*` Images and videos
- `/styles.css`, `/admin.css`, `/invitation.css`, `/portal.css`, `/scan.css`, `/monitor.css`
- `/spinner.js`, `/admin.js`, `/portal.js`, `/scan.js`, `/monitor.js`

## Authentication
- Basic Auth for `/admin` and `/spinner`:
  - Username: `admin`
  - Password: `Dptaspen36`
- API requests require `X-API-Key: Dptaspen@25!`
- WebSocket requires `key` query parameter: `/ws/attendance?key=Dptaspen@25!`

## API Endpoints
All endpoints under `/api` require `X-API-Key`.

### Employees
- `GET /api/employees/all`
- `GET /api/employees/present`
- `POST /api/employees/mark_present`
- `DELETE /api/employees/present`
- `GET /api/employees/export`

### Guests
- `GET /api/guests`
- `POST /api/guests/mark_present`
- `DELETE /api/guests/present`

### Winners
- `GET /api/winners`
- `POST /api/winners`
- `DELETE /api/winners`
- `GET /api/winners/export`

### Invitations
- `GET /api/invitations/lookup?nip=...`

### Scans (Realtime Monitoring)
- `POST /api/scans` (payload: `name`, `scanner_id`)
- `GET /api/scans/recent`

### WebSocket
- `GET /ws/attendance` (broadcasts scan events)

## Scanner Configuration
Open the scanner page with a scanner ID:
- `/scan?scanner=1`
- `/scan?scanner=2`
- `/scan?scanner=3`

If the query parameter is missing or invalid, the page will prompt for a valid scanner ID (1-3).

## Database
The app uses MySQL. The main tables are:
- `employees` (includes `nip`)
- `attendances` (employee and guest attendance)
- `winners`
- `scan_events`

Migrations are located in the `migrations` directory and are managed by Goose.

## Setup
1. Prepare MySQL and create a database.
2. Configure the connection string in `.env`:

```bash
GOOSE_DBSTRING="user:password@tcp(host:3306)/database?parseTime=true"
```

3. Run migrations:

```bash
goose -dir migrations mysql "$GOOSE_DBSTRING" up
```

4. Start the server:

```bash
go run cmd/web/main.go
```

5. Open the pages:
- `http://localhost:8080/portal`
- `http://localhost:8080/scan?scanner=1`
- `http://localhost:8080/monitor`
- `http://localhost:8080/admin`
- `http://localhost:8080/spinner`

## Notes
- The invitation page generates QR codes using a third-party service.
- API key and Basic Auth credentials are currently hardcoded in the server and frontend assets.
- Scan events are deduplicated by name on the server; only the first scan is stored.

## Project Layout
- `cmd/web` Server entrypoint
- `internal/adapter/http/handler` HTTP handlers and WebSocket hub
- `internal/model` Data models
- `internal/repository` Database access
- `internal/services` Business logic
- `internal/templates/assets` Templates, CSS, JS, and static assets
- `migrations` Goose SQL migrations
