#!/usr/bin/env bash
set -euo pipefail
docker container stop hotdeals_postgres-dev
docker container stop hotdeals_scraper-dev
docker container stop hotdeals_mailer-dev
docker-compose -f /drone/src/docker-compose.yml -f /drone/src/docker-compose.prod.yml up --detach --build
