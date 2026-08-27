.PHONY: demo demo-down demo-shell demo-status

# Override e.g. `make demo SALT_VERSION=3006.26 PYTHON_VERSION=3.10`
SALT_VERSION ?= 3008.2
PYTHON_VERSION ?= 3.13

demo:
	SALT_VERSION=$(SALT_VERSION) PYTHON_VERSION=$(PYTHON_VERSION) \
		docker compose -f e2e_test/docker-compose.yaml -f e2e_test/docker-compose.demo.yaml up --build

demo-down:
	docker compose -f e2e_test/docker-compose.yaml -f e2e_test/docker-compose.demo.yaml down

demo-shell:
	docker compose -f e2e_test/docker-compose.yaml -f e2e_test/docker-compose.demo.yaml exec salt_master bash

demo-status:
	docker compose -f e2e_test/docker-compose.yaml -f e2e_test/docker-compose.demo.yaml ps
