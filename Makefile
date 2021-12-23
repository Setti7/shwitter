
dev: up
	cd src && air

build:
	go build cmd/shwitter/shwitter.go

up:
	docker-compose up -d

csql: up
	docker-compose exec cassandra cqlsh

create-migration:
	migrate create -ext cql -dir ./migrations $(name)
