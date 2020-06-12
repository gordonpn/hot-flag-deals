#!/bin/sh
echo "$DOCKER_TOKEN" | docker login -u gordonpn --password-stdin
docker run --rm --privileged multiarch/qemu-user-static --reset -p yes
docker buildx rm builder || true
docker buildx create --name builder --driver docker-container --use
docker buildx inspect --bootstrap
cd /drone/src/mailer || exit 1
docker buildx build -t gordonpn/hotdeals-mailer:latest --platform linux/amd64,linux/arm64 --push .
cd /drone/src/scraper || exit 1
docker buildx build -t gordonpn/hotdeals-scraper:latest --platform linux/amd64,linux/arm64 --push .
# todo buld for backend, frontend and proxy
# todo tag appropriately stable vs latest
