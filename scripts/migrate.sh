#!/bin/bash

# Migration script for Railway deployment
# This script runs all migration files in order

set -e

echo "🚀 Starting database migration..."

# Check if DATABASE_URL is set
if [ -z "$DATABASE_URL" ]; then
    echo "❌ ERROR: DATABASE_URL environment variable is not set"
    exit 1
fi

echo "📊 Database URL: ${DATABASE_URL%%@*}@****"

# Run migrations in order
echo "📝 Running migrations..."

# Run init-db.sql (combined schema)
if [ -f "./scripts/init-db.sql" ]; then
    echo "  → Running init-db.sql..."
    psql "$DATABASE_URL" -f ./scripts/init-db.sql
else
    echo "⚠️  init-db.sql not found, running individual migrations..."

    # Run migrations from migrations folder
    for file in ./migrations/*.up.sql; do
        if [ -f "$file" ]; then
            echo "  → Running $(basename $file)..."
            psql "$DATABASE_URL" -f "$file"
        fi
    done
fi

# Run seed data
if [ -f "./scripts/seed-data.sql" ]; then
    echo "🌱 Seeding demo data..."
    psql "$DATABASE_URL" -f ./scripts/seed-data.sql
fi

echo "✅ Migration completed successfully!"
