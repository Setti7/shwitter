cass_driver = cassandra://127.0.0.1:9042/shwitter

migrate_cmd = migrate -database $(cass_driver) -path ./migrations

dev: up
	cd src && air

up:
	docker-compose up -d

csql: up
	docker-compose exec cass_cluster cqlsh

migrate:
	echo "$(migrate_cmd)"

create-migration:
	$(migrate_cmd) migrate create $(name)

migrate-up:
	$(migrate_cmd) up

migrate-down:
	$(migrate_cmd) down

clear-db:
	$(migrate_cmd) down
	$(migrate_cmd) up