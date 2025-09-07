# Database Setup Guide

## Prerequisites

1. PostgreSQL 12+ installed and running
2. Database created for the application

## Environment Variables

Create a `.env` file in the project root with the following variables:

```env
# Database Configuration
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_NAME=cupid
DB_PASSWORD=your_password
DB_URL=postgres://postgres:your_password@localhost:5432/cupid?sslmode=disable

# Cupid API Configuration
CUPID_API_BASE_URL=https://api.cupid.com
CUPID_API_VERSION=v1
CUPID_API_KEY=your_api_key_here

# Server Configuration
PORT=8080
ENV=development
```

## Database Setup

1. **Create the database:**
   ```bash
   createdb cupid
   ```

2. **Run migrations:**
   ```bash
   make migrate-up
   ```

3. **Check migration status:**
   ```bash
   make migrate-status
   ```

## Available Commands

- `make migrate-up` - Run all pending migrations
- `make migrate-down` - Rollback the last migration
- `make migrate-status` - Check migration status
- `make migrate-create NAME=create_table` - Create a new migration
- `make migrate-reset` - Reset database (WARNING: deletes all data)
- `make db-setup` - Create database if it doesn't exist

## Database Schema

The application uses a hybrid approach:

### Normalized Tables (for fast queries):
- `properties` - Main property data with essential fields
- `reviews` - Property reviews with searchable fields
- `translations` - Property translations by language

### JSONB Tables (for complex data):
- `property_details` - Complex nested data stored as JSONB

## Testing the Setup

1. **Run the fetch command to populate data:**
   ```bash
   make fetch-data
   ```

2. **Start the API server:**
   ```bash
   make run
   ```

3. **Check the API documentation:**
   Visit `http://localhost:8080/docs` for Swagger documentation
