#!/bin/bash

# Cupid API Docker Startup Script
# This script sets up and starts the Cupid API with PostgreSQL

set -e

echo "🚀 Starting Cupid API with Docker Compose..."

# Check if docker-compose is available
if ! command -v docker-compose &> /dev/null; then
    echo "❌ docker-compose is not installed. Please install Docker Compose first."
    exit 1
fi

# Check if docker is running
if ! docker info &> /dev/null; then
    echo "❌ Docker is not running. Please start Docker first."
    exit 1
fi

# Create docker.env if it doesn't exist
if [ ! -f "docker.env" ]; then
    echo "📝 Creating docker.env from integration.env.example..."
    cp integration.env.example docker.env
    echo "⚠️  Please edit docker.env and set your CUPID_API_KEY before running again."
    echo "   You can also modify other settings as needed."
    exit 1
fi

# Check if CUPID_API_KEY is set
if ! grep -q "CUPID_API_KEY=" docker.env || grep -q "CUPID_API_KEY=your_api_key_here" docker.env; then
    echo "❌ Please set your CUPID_API_KEY in docker.env file"
    echo "   Edit docker.env and replace 'your_api_key_here' with your actual API key"
    exit 1
fi

echo "🔧 Building and starting services..."

# Build and start the services
docker-compose up --build -d

echo "⏳ Waiting for services to be ready..."

# Wait for PostgreSQL to be ready
echo "   Waiting for PostgreSQL..."
timeout 60 bash -c 'until docker-compose exec postgres pg_isready -U root -d cupid; do sleep 2; done'

# Wait for API to be ready
echo "   Waiting for API..."
timeout 60 bash -c 'until curl -f http://localhost:8080/api/v1/health 2>/dev/null; do sleep 2; done'

echo "✅ Services are ready!"

# Run database migrations
echo "🗄️  Running database migrations..."
docker-compose exec api sh -c "cd /app && ./api" &
API_PID=$!

# Wait a bit for the API to start
sleep 10

# Run migrations (you might need to implement this based on your migration setup)
echo "   Migrations will be handled by the application startup"

# Fetch initial data
echo "📥 Fetching initial data from Cupid API..."
docker-compose run --rm fetch-data

echo ""
echo "🎉 Cupid API is now running!"
echo ""
echo "📊 Services:"
echo "   • API: http://localhost:8080"
echo "   • PostgreSQL: localhost:5432"
echo "   • Database: cupid"
echo ""
echo "📚 API Documentation:"
echo "   • Swagger UI: http://localhost:8080/swagger/index.html"
echo ""
echo "🔧 Management Commands:"
echo "   • View logs: docker-compose logs -f"
echo "   • Stop services: docker-compose down"
echo "   • Restart: docker-compose restart"
echo "   • Fetch new data: docker-compose run --rm fetch-data"
echo ""
echo "📝 To view the API logs:"
echo "   docker-compose logs -f api"
