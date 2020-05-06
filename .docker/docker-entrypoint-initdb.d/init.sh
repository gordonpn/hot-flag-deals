#!/usr/bin/env bash
set -euo pipefail

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
		CREATE USER ${POSTGRES_NONROOT_USER} WITH PASSWORD '${POSTGRES_NONROOT_PASSWORD}';
		CREATE DATABASE deals OWNER ${POSTGRES_NONROOT_USER} ENCODING 'UTF8';
    GRANT ALL PRIVILEGES ON DATABASE deals TO ${POSTGRES_NONROOT_USER};
EOSQL
