SHELL := /bin/bash

AL_TMPL  := prometheus/alertmanager.tmpl.yml
AL_OUT   := prometheus/generated/alertmanager.yml

COMPOSE  := docker compose

.PHONY: render build up down restart logs prune

render:
        @mkdir -p $(dir $(AL_OUT))
        @set -a; source ./.env; set +a; \
         envsubst < $(AL_TMPL) > $(AL_OUT)

up: render
        $(COMPOSE) up -d

down:
        $(COMPOSE) down

restart: down up

build:
        $(COMPOSE) build

logs:
        $(COMPOSE) logs -f --tail=100

prune:
        docker system prune -f --volumes



