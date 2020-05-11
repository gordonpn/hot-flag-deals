#!/usr/bin/env bash
set -euo pipefail
docker container stop hotdeals_postgres || true
docker container stop hotdeals_scraper || true
docker container stop hotdeals_mailer || true
docker container rm hotdeals_postgres || true
docker container rm hotdeals_scraper || true
docker container rm hotdeals_mailer || true
docker-compose -f /drone/src/docker-compose.yml -f /drone/src/docker-compose.dev.yml up --detach --build
