all: install run help

install: build_go build_admin

run:
	docker-compose up -d
	docker exec -ti -d nebuleuse-go bash -c "./Nebuleuse"

build_go:
	docker-compose up -d nebuleuse-go nebuleuse-database
	docker exec -ti nebuleuse-go bash -c "cp -n .config.docker .config"
	docker exec -ti nebuleuse-go bash -c "go get ."
	docker exec -ti nebuleuse-go bash -c "go build"

build_admin:
	docker-compose up -d nebuleuse-node
	docker exec -ti nebuleuse-node bash -c "npm install"
	docker exec -ti nebuleuse-node bash -c "npm run bower-install"
	docker exec -ti nebuleuse-node bash -c "npm run build"

rebuild_admin:
	docker-compose up -d nebuleuse-node
	docker exec -ti nebuleuse-node bash -c "npm run build"

bash_go:
	docker-compose up -d nebuleuse-go
	docker exec -ti nebuleuse-go bash

bash_node:
	docker-compose up -d nebuleuse-node
	docker exec -ti nebuleuse-node bash

logs:
	docker-compose logs -ft nebuleuse-go

help:
	#
	# Nebuleuse
	#
	# Nebuleuse API:    http://0.0.0.0:12080/
	# Nebuleuse Admin:  http://0.0.0.0:12080/admin/dist/ (user: test / pass: test)
	#
