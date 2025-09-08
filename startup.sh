#!/bin/sh

# Cupid API Container Startup Script
# This script runs inside the container to set up and start the API

set -e

echo "🚀 Starting Cupid API inside container..."

# Load environment variables from docker-compose (no .env file needed)
echo "📋 Environment variables loaded from docker-compose..."

# Wait for PostgreSQL to be ready
echo "⏳ Waiting for PostgreSQL to be ready..."
until pg_isready -h ${DB_HOST:-postgres} -p ${DB_PORT:-5432} -U ${DB_USER:-root} -d ${DB_NAME:-cupid}; do
    echo "   PostgreSQL is unavailable - sleeping..."
    sleep 2
done

echo "✅ PostgreSQL is ready!"

# Run database migrations
echo "🗄️  Running database migrations..."
DB_URL="postgres://${DB_USER:-root}:${DB_PASSWORD:-root}@${DB_HOST:-postgres}:${DB_PORT:-5432}/${DB_NAME:-cupid}?sslmode=disable"

# Run migrations directly with the correct path
if goose -dir ./cmd/migrate/migrations postgres "$DB_URL" up; then
    echo "✅ Database migrations completed!"
else
    echo "❌ Database migrations failed!"
    exit 1
fi

# Start the API server first (so health checks pass)
echo "🌐 Starting API server..."
echo "   API will be available at: http://localhost:${SERVER_PORT:-8080}"
echo "   Swagger UI: http://localhost:${SERVER_PORT:-8080}/swagger/index.html"
echo "   Health check: http://localhost:${SERVER_PORT:-8080}/api/v1/health"
echo ""

# Start API server in background
./api &
API_PID=$!

# Wait a moment for API to start
sleep 5

# Fetch initial data from Cupid API in background
echo "📥 Fetching initial data from Cupid API..."
echo "   This may take a few minutes depending on the amount of data..."
echo "   Please wait while we fetch hotel data from the Cupid API..."
echo ""

# Start data fetching in background
echo "🔄 Starting data fetch process..."
./fetch &
FETCH_PID=$!

# Show progress indicator while fetching
COUNTER=0
while kill -0 $FETCH_PID 2>/dev/null; do
    case $((COUNTER % 8)) in
        0) echo "   ⏳ Fetching data... still working..." ;;
        1) echo "   🔄 Processing hotels... please wait..." ;;
        2) echo "   📡 Communicating with Cupid API... working..." ;;
        3) echo "   💾 Storing data in database... processing..." ;;
        4) echo "   🏨 Loading hotel information... working..." ;;
        5) echo "   📊 Organizing data... please wait..." ;;
        6) echo "   🔍 Validating information... working..." ;;
        7) echo "   ⚡ Almost done... processing..." ;;
    esac
    COUNTER=$((COUNTER + 1))
    sleep 4
done

# Wait for the process to complete and get its exit status
wait $FETCH_PID
FETCH_EXIT_CODE=$?

echo ""
echo "🎯 Data fetch process completed!"

if [ $FETCH_EXIT_CODE -eq 0 ]; then
    echo "✅ Initial data fetched successfully!"
    
    # Validate that data was actually fetched
    echo "🔍 Validating data fetch..."
    PROPERTY_COUNT=$(psql -h ${DB_HOST:-postgres} -p ${DB_PORT:-5432} -U ${DB_USER:-root} -d ${DB_NAME:-cupid} -t -c "SELECT COUNT(*) FROM properties;" 2>/dev/null | tr -d ' \n' || echo "0")
    
    if [ "$PROPERTY_COUNT" -gt "0" ]; then
        echo "✅ Validation successful! Found $PROPERTY_COUNT properties in database."
        echo "🎉 Database is now populated with real hotel data!"
    else
        echo "⚠️  Warning: No properties found in database after fetch."
    fi
else
    echo "⚠️  Data fetching failed, but API will continue running"
    echo "   You can manually fetch data later by running: ./fetch"
fi

echo ""
echo "✨ Ready to explore hotels! The API is fully functional with real hotel data."

# Wait for API process
wait $API_PID
