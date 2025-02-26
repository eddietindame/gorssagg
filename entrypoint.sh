#!/bin/bash
set -e

MAX_ATTEMPTS=10
ATTEMPT=1

# Wait for the database to be ready
until PGPASSWORD="$POSTGRES_PASSWORD" psql -h postgres_db -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c '\q' >/dev/null 2>&1; do
  if [ "$ATTEMPT" -ge "$MAX_ATTEMPTS" ]; then
    echo "Postgres is not ready after $MAX_ATTEMPTS attempts. Failing..."
    exit 1
  fi
  echo "Attempt $ATTEMPT: Waiting for Postgres..."
  sleep 2
  ATTEMPT=$((ATTEMPT + 1))
done

echo "Database is ready!"

# Run Goose migrations
echo "Running migrations..."
if ! goose -dir /usr/app/sql/schema postgres "postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@postgres_db:5432/$POSTGRES_DB?sslmode=disable" up; then
  echo "Database migrations failed!"
  exit 1
fi

# Start the app
if [ $# -eq 0 ]; then
  echo "No command provided to run. Exiting..."
  exit 1
fi

echo "Starting the app..."
exec "$@"
