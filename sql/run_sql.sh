#!/bin/bash

# Database configuration
DB_NAME="gochat"
DB_USER="gochat"
DB_PASSWORD="password"

if [ -z "$1" ]; then
  echo "Usage: $0 <sql_file>"
  exit 1
fi

PGPASSWORD="$DB_PASSWORD" psql -h localhost -p 5432 -U "$DB_USER" -d "$DB_NAME" -f "$1"
