#!/bin/bash

# Cupid API Docker Startup Script
# This script sets up and starts the Cupid API with PostgreSQL

set -e

echo "🚀 Starting Cupid API with Docker Compose..."

# Check if docker-compose or docker compose is available
DOCKER_COMPOSE_CMD=""
if command -v docker-compose &> /dev/null; then
    DOCKER_COMPOSE_CMD="docker-compose"
elif docker compose version &> /dev/null; then
    DOCKER_COMPOSE_CMD="docker compose"
else
    echo "❌ Neither 'docker-compose' nor 'docker compose' is available. Please install Docker Compose first."
    exit 1
fi

echo "✅ Using Docker Compose command: $DOCKER_COMPOSE_CMD"

# Load environment variables from docker.env
if [ -f "docker.env" ]; then
    echo "📋 Loading environment variables from docker.env..."
    export $(grep -v '^#' docker.env | grep -v '^$' | xargs)
fi

# Check if docker is running
if ! docker info &> /dev/null; then
    echo "❌ Docker is not running. Please start Docker first."
    exit 1
fi

# Create docker.env if it doesn't exist
if [ ! -f "docker.env" ]; then
    echo "📝 Creating docker.env from docker.env.example..."
    cp docker.env.example docker.env
    echo "⚠️  Please edit docker.env and set your CUPID_API_KEY before running again."
    echo "   You can also modify other settings as needed."
    exit 1
fi

# Check if required environment variables are set in docker.env
REQUIRED_VARS=(
    "CUPID_API_KEY"
    "DB_HOST"
    "DB_PORT"
    "DB_USER"
    "DB_PASSWORD"
    "DB_NAME"
    "DB_DRIVER"
    "SERVER_PORT"
    "GO_ENV"
    "LOG_LEVEL"
)

MISSING_VARS=()
DEFAULT_PLACEHOLDERS=(
    ["CUPID_API_KEY"]="your_api_key_here"
    ["DB_HOST"]="your_database_host"
    ["DB_USER"]="your_database_user"
    ["DB_PASSWORD"]="your_database_password"
    ["DB_NAME"]="your_database_name"
)

for VAR in "${REQUIRED_VARS[@]}"; do
    # Check if variable is present
    if ! grep -q "^${VAR}=" docker.env; then
        MISSING_VARS+=("$VAR (missing)")
        continue
    fi
    # Check for placeholder value if defined
    PLACEHOLDER="${DEFAULT_PLACEHOLDERS[$VAR]}"
    if [ -n "$PLACEHOLDER" ]; then
        if grep -q "^${VAR}=${PLACEHOLDER}" docker.env; then
            MISSING_VARS+=("$VAR (placeholder: $PLACEHOLDER)")
        fi
    fi
done

if [ ${#MISSING_VARS[@]} -ne 0 ]; then
    echo "❌ The following required environment variables are missing or have placeholder values in docker.env:"
    for VAR in "${MISSING_VARS[@]}"; do
        echo "   - $VAR"
    done
    echo ""
    echo "   Please edit docker.env and set the correct values before running again."
    exit 1
fi

echo "🔧 Building and starting services..."

# Build and start the services
$DOCKER_COMPOSE_CMD up --build -d

echo "⏳ Waiting for services to be ready..."

# Wait for services to be healthy
echo "⏳ Waiting for services to be healthy..."
echo "   Checking PostgreSQL health..."

# Wait for PostgreSQL to be healthy
for i in {1..30}; do
    if $DOCKER_COMPOSE_CMD exec -T postgres pg_isready -U ${DB_USER:-root} -d ${DB_NAME:-cupid} > /dev/null 2>&1; then
        echo "   ✅ PostgreSQL is healthy!"
        break
    fi
    echo "   ⏳ PostgreSQL not ready yet... (attempt $i/30)"
    sleep 2
done

echo "   Checking API health..."
# Wait for API to be healthy
for i in {1..30}; do
    if curl -s http://localhost:${SERVER_PORT:-8080}/api/v1/health > /dev/null 2>&1; then
        echo "   ✅ API is healthy!"
        break
    fi
    echo "   ⏳ API not ready yet... (attempt $i/30)"
    sleep 2
done

echo "✅ All services are ready!"

# Data fetching is now handled automatically inside the API container
echo "📝 Data fetching is handled automatically by the API container startup"
echo "   The API container will fetch hotel data from Cupid API during startup"
echo "   This process may take several minutes - please be patient!"
echo "   You can monitor progress with: $DOCKER_COMPOSE_CMD logs -f api"
echo ""
echo "   💡 Tip: The container will show progress messages during data fetching"

echo ""
echo "🎉 Cupid API is now running!"
echo ""
echo "📊 Services:"
echo "   • API: http://localhost:${SERVER_PORT:-8080}"
echo "   • PostgreSQL: localhost:5432"
echo "   • Database: ${DB_NAME:-cupid}"
echo ""
echo "📚 API Documentation:"
echo "   • Swagger UI: http://localhost:${SERVER_PORT:-8080}/docs/index.html"
echo "   • Health Check: http://localhost:${SERVER_PORT:-8080}/api/v1/health"
echo ""
echo "🧪 Quick Test:"
echo "   • List properties: curl http://localhost:${SERVER_PORT:-8080}/api/v1/properties"
echo "   • Search hotels: curl 'http://localhost:${SERVER_PORT:-8080}/api/v1/search?city=Paris'"
echo ""
echo "🔧 Management Commands:"
echo "   • View logs: $DOCKER_COMPOSE_CMD logs -f"
echo "   • Stop services: $DOCKER_COMPOSE_CMD down"
echo "   • Restart: $DOCKER_COMPOSE_CMD restart"
echo "   • Fetch new data: $DOCKER_COMPOSE_CMD exec api ./fetch"
echo ""
echo "📝 To view the API logs:"
echo "   $DOCKER_COMPOSE_CMD logs -f api"
echo ""
echo "✨ Ready to explore hotels! The API is fully functional with real hotel data."
