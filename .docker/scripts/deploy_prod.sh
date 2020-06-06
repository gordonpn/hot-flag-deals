#!/bin/sh
echo "Not implemented"
exit 1
docker container stop hotdeals_postgres-dev || true
docker container stop hotdeals_scraper-dev || true
docker container stop hotdeals_mailer-dev || true
docker container rm hotdeals_postgres-dev || true
docker container rm hotdeals_scraper-dev || true
docker container rm hotdeals_mailer-dev || true
docker-compose -f /drone/src/docker-compose.yml -f /drone/src/docker-compose.prod.yml up --detach --build
