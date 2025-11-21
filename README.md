# LocalHaven CMS Site

Website and survey application for LocalHaven CMS.

## Features

- Modern, responsive web interface
- Multi-step survey form
- Analytics dashboard
- Admin interface
- Local-first architecture

## Technology Stack

- Frontend: Astro + Svelte
- Backend: Go
- Database: SQLite
- Container: Docker

## Development Setup

### Prerequisites

- Docker and Docker Compose
- Node.js 20+ (for local development without Docker)
- Go 1.21+ (for local development without Docker)

### Running with Docker (Development)

1. Clone the repository:

```bash
git clone https://github.com/yourusername/LocalHavenCMS-Site.git
cd LocalHavenCMS-Site
```

2. Start the development environment:

```bash
docker compose -f docker-compose.dev.yml up
```

The development environment will be available at:

- Frontend: http://localhost:3000 (with hot reload)
- API: http://localhost:8090

To stop the development environment:

```bash
docker compose -f docker-compose.dev.yml down
```

### Running Locally (Without Docker)

1. Install frontend dependencies:

```bash
cd web
npm install
```

2. Start the frontend:

```bash
npm run dev
```

3. Install backend dependencies:

```bash
cd backend
go mod download
```

4. Start the backend:

```bash
go run main.go
```

### Environment Variables

The following environment variables can be configured:

```bash
# Required
JWT_SECRET=your_jwt_secret
ADMIN_USERNAME=admin
ADMIN_PASSWORD=your_secure_password

# Optional
PORT=8090
ENVIRONMENT=development
ALLOWED_ORIGINS=http://localhost:3000
TRUSTED_PROXIES=172.16.0.0/12,192.168.0.0/16,10.0.0.0/8,127.0.0.1
```

## Development

- Frontend code is in the `web` directory
- Backend code is in the `backend` directory
- Database files are stored in the `data` directory

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## License

MIT License - See LICENSE file for details
