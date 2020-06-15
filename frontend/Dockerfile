FROM node:14-alpine AS builder
RUN mkdir -p /usr/src/app && chown -R node:node /usr/src/app
WORKDIR /usr/src/app
ENV PATH /usr/src/app/node_modules/.bin:$PATH
COPY --chown=node:node package*.json ./
USER node
RUN npm ci --only=production
RUN npm cache clean --force --loglevel=error
COPY --chown=node:node . .
RUN npm run export

FROM caddy:2-alpine
RUN apk add --no-cache \
  curl \
  && rm -rf /var/cache/apk/*
COPY --from=builder /usr/src/app/out /usr/share/caddy
COPY Caddyfile /etc/caddy/Caddyfile
EXPOSE 3000
