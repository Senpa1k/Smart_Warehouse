#!/bin/sh

echo "Running migrations..."
./migrator || echo "Migrations completed with warnings (tables may already exist)"

echo "Starting application..."
exec ./main
