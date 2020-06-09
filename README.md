# Hot Deals

## Description

Software as a service scraping the [Hot Deals forums](https://forums.redflagdeals.com/hot-deals-f9/) and delivering newsletter emails on trending deals from the Forum.

---
[![Build Status](https://drone.gordon-pn.com/api/badges/gordonpn/hot-flag-deals/status.svg)](https://drone.gordon-pn.com/gordonpn/hot-flag-deals)
![Healthchecks.io](https://healthchecks.io/badge/85143171-fec6-42e7-b3d8-bc7f499f5d0d/r5YwfK0Y.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/gordonpn/hot-flag-deals)](https://goreportcard.com/report/github.com/gordonpn/hot-flag-deals)
![Last commit on develop](https://badgen.net/github/last-commit/gordonpn/hot-flag-deals/develop)
![License](https://badgen.net/github/license/gordonpn/hot-flag-deals)

[![Buy Me A Coffee](https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png)](https://www.buymeacoffee.com/gordonpn)

## Objective / Motivation

Red Flag Deals does aggregate deals on their front page, **but** the Hot Deals Forums are community driven and crowdsourced. This is where the purpose of my project comes into play, this project scrapes the Hot Deals Forums several times per day.

With this project, I saved myself the chore of checking the (messy) forum a few times a day while still being aware of the good deals posted by the community.

## Screenshot

<a href="./docs/newsletter.png"><img src="./docs/newsletter.png" height="900"></a>

Base template design by [@tiffzeng](https://github.com/tiffzeng)

## How It Works

The project is split into several services:

* Scraper

  Scrapes the forums periodically and depends on the database to persist the data.

* Mailer

  Sends the daily newsletter email.

* Backend

  RESTful API, takes care of sending data to the frontend and managing subscribing and unsubscribing.

* Frontend

  Displays the deals nicely, and has the subscription form.

* Proxy

  Forwards requests to the right service.

## Built with / technologies

* Go programming language
* SendGrid
* Docker & Docker Swarm
* PostgreSQL
* Drone CI
* Nginx
* Next.js
* React.js
* Material-UI

## Features

* Automated scraping of the Forum
* Insert and update a database
* Automated newsletter email
* Algorithm to filter the junk
* Website to browse the deals and subscribe to the emails

## Getting started

### Prerequisites

* Go(lang) v1.14+
* Docker-compose v1.25.5+
* Docker v19.03+
* Node.js v14+

### Configuration

| Variable                  | Description                                                                                                |
|---------------------------|------------------------------------------------------------------------------------------------------------|
| SENDGRID_API_KEY          | SendGrid API key to use the service                                                                        |
| SENDGRID_TEMPLATE         | SendGrid requires you to upload an html template to use for dynamic emails, this is the ID of the template |
| POSTGRES_NONROOT_DB       | Database for the appuser                                                                                   |
| POSTGRES_DB               | Root database                                                                                              |
| POSTGRES_USER             | Root username to initialize PostgreSQL                                                                     |
| POSTGRES_PASSWORD         | Root password to initialize PostgreSQL                                                                     |
| POSTGRES_NONROOT_PASSWORD | The user password to use the database, we don't want to use the root user for these tasks                  |
| POSTGRES_NONROOT_USER     | The user username to use the database                                                                      |
| DEV                       | Set this variable to anything to run as dev (testing purposes)                                             |
| SCRAPER_HC_UUID           | UUID for Healthchecks.io                                                                                   |
| MAILER_HC_UUID            | UUID for Healthchecks.io                                                                                   |

### Installing

Build individual Go services into a binary

`go build -o main .`

### Usage

Bring the services up using docker-compose.

## Running the tests

`go test`

## Roadmap / Todo

Check out the [open issues](https://github.com/gordonpn/hot-flag-deals/issues?q=is%3Aissue+is%3Aopen+sort%3Aupdated-desc) for ideas and features I have planned!

## Support

You may open an issue for discussion.

## Authors

* Myself [@gordonpn](https://github.com/gordonpn)

## License

[MIT License](./LICENSE)
