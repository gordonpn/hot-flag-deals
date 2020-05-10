#!/usr/bin/env bash
set -euo pipefail
docker container stop hotdeals_postgres
docker container stop hotdeals_scraper
docker container stop hotdeals_mailer
docker-compose -f /drone/src/docker-compose.yml -f /drone/src/docker-compose.dev.yml up --detach --build
