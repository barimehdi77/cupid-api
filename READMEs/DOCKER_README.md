# Cupid API Docker Setup

This document explains how to run the Cupid API using Docker and Docker Compose.

## Prerequisites

- Docker and Docker Compose installed
- A valid Cupid API key

## Quick Start

### 1. Configure Environment

Copy the example environment file and set your API key:

```bash
cp docker.env.example docker.env
```

Edit `docker.env` and set your `CUPID_API_KEY`:

```bash
# Edit docker.env
CUPID_API_KEY=your_actual_api_key_here
```

### 2. Start the Services

Use the provided startup script:

```bash
./docker-startup.sh
```

Or manually with Docker Compose:

```bash
# Build and start all services
docker-compose up --build -d

# Wait for services to be ready, then fetch data
docker-compose run --rm fetch-data
```

## Services

The Docker setup includes the following services:

### PostgreSQL Database

- **Container**: `cupid-postgres`
- **Port**: `5432`
- **Database**: `cupid`
- **User**: `root`
- **Password**: `root`

### Cupid API

- **Container**: `cupid-api`
- **Port**: `8080`
- **Health Check**: `http://localhost:8080/health`
- **Swagger UI**: `http://localhost:8080/docs/index.html`

### Data Fetcher

- **Container**: `cupid-fetch-data` (runs on demand)
- **Purpose**: Fetches hotel data from Cupid API and stores it in PostgreSQL

## Environment Variables

The following environment variables can be configured in `docker.env`:

| Variable | Default | Description |
|----------|---------|-------------|
| `CUPID_API_KEY` | - | **Required**: Your Cupid API key |
| `CUPID_API_BASE_URL` | `https://content-api.cupid.travel` | Cupid API base URL |
| `CUPID_API_VERSION` | `v3.0` | Cupid API version |
| `DB_NAME` | `cupid` | PostgreSQL database name |
| `DB_USER` | `root` | PostgreSQL username |
| `DB_PASSWORD` | `root` | PostgreSQL password |
| `SERVER_PORT` | `8080` | API server port |
| `GO_ENV` | `production` | Go environment |
| `LOG_LEVEL` | `info` | Logging level |

## Management Commands

### Start Services

```bash
docker-compose up -d
```

### Stop Services

```bash
docker-compose down
```

### View Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f api
docker-compose logs -f postgres
```

### Fetch New Data

```bash
docker-compose run --rm fetch-data
```

### Restart Services

```bash
docker-compose restart
```

### Rebuild and Start

```bash
docker-compose up --build -d
```

### Access Database

```bash
# Connect to PostgreSQL
docker-compose exec postgres psql -U root -d cupid

# Or use external client
psql -h localhost -p 5432 -U root -d cupid
```

### Access API Container

```bash
docker-compose exec api sh
```

## Data Persistence

- PostgreSQL data is persisted in a Docker volume named `postgres_data`
- Data survives container restarts and rebuilds
- To reset the database, remove the volume: `docker-compose down -v`

## Health Checks

Both services include health checks:

- **PostgreSQL**: Checks if the database is ready to accept connections
- **API**: Checks if the HTTP health endpoint responds

## Troubleshooting

### Services Won't Start

1. Check if Docker is running: `docker info`
2. Check logs: `docker-compose logs`
3. Ensure port 8080 and 5432 are not in use

### API Key Issues

1. Verify your API key in `docker.env`
2. Check if the key is valid by testing the Cupid API directly

### Database Connection Issues

1. Wait for PostgreSQL to be ready (health check)
2. Check database logs: `docker-compose logs postgres`
3. Verify database credentials in `docker.env`

### Data Fetching Issues

1. Check API logs: `docker-compose logs api`
2. Verify Cupid API connectivity
3. Check fetch-data logs: `docker-compose logs fetch-data`

## Development

### Local Development with Docker

```bash
# Start only the database
docker-compose up postgres -d

# Run the API locally
make dev
```

### Testing

```bash
# Run tests in container
docker-compose exec api go test ./...

# Run integration tests
docker-compose exec api go test -tags=integration ./...
```

## Production Considerations

For production deployment:

1. **Security**:
   - Change default database credentials
   - Use environment-specific API keys
   - Enable SSL for database connections

2. **Performance**:
   - Adjust PostgreSQL configuration
   - Set appropriate resource limits
   - Use a reverse proxy (nginx)

3. **Monitoring**:
   - Add monitoring and logging solutions
   - Set up health check endpoints
   - Monitor resource usage

4. **Backup**:
   - Implement database backup strategy
   - Regular data synchronization with Cupid API

## API Endpoints

Once running, the API provides:

- `GET /health` - Health check
- `GET /swagger/index.html` - API documentation
- `GET /api/v1/properties` - List properties
- `GET /api/v1/properties/{id}` - Get property details
- `POST /api/v1/admin/sync` - Trigger data sync
- `GET /api/v1/admin/sync/status` - Check sync status

## Support

For issues or questions:

1. Check the logs: `docker-compose logs`
2. Verify configuration in `docker.env`
3. Test individual components
4. Check the main README.md for additional information
