cass_driver = cassandra://127.0.0.1:9042/shwitter

journey_cmd = journey --url $(cass_driver) --path ./migrations

dev: up
	cd src && air

up:
	docker-compose up -d

csql: up
	docker-compose exec cass_cluster cqlsh

create-keyspace: up
	echo "CREATE KEYSPACE IF NOT EXISTS shwitter WITH replication = {'class': 'SimpleStrategy', 'replication_factor' : 1};" | docker-compose exec -it cass_cluster cqlsh

j:
	echo "$(journey_cmd)"

create-migration:
	$(journey_cmd) migrate create $(name)

migrate-up:
	$(journey_cmd) migrate up

migrate-down:
	$(journey_cmd) migrate down

clear-db:
	$(journey_cmd) migrate down
	$(journey_cmd) migrate up