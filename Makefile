DCOMPOSE = docker-compose -f scripts/docker-compose.yml

.PHONY: up down restart logs

up:
	$(DCOMPOSE) up --build

down:
	$(DCOMPOSE) down

restart:
	$(DCOMPOSE) down && $(DCOMPOSE) up --build

logs:
	$(DCOMPOSE) logs -f
