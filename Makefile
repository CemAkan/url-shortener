SHELL := /bin/bash

COMPOSE  := docker compose

.PHONY: build up down restart logs prune

up:
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



