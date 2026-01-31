There are multiple ways to run this project:

**Option 1: Docker (Full Stack)**
```bash
docker-compose up
```

**Option 2: Development Mode (Individual Services)**

For the **web frontend**:
```bash
cd web && pnpm install && pnpm dev
```

For the **backend** (Go):
```bash
cd backend && go run .
```

Which would you like me to run?
- **Docker Compose** - runs both backend and frontend together
- **Frontend only** - just the Astro dev server
- **Backend only** - just the Go API server
