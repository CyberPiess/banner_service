all:
	docker-compose -f build/docker-compose.yaml up -d cache
	docker-compose -f build/docker-compose.yaml up -d db
	docker exec -it local_pgsql createdb --username=myuser --owner=myuser banner_service
	migrate -path db/migration -database "postgresql://myuser:mypassword@localhost:5432/banner_service?sslmode=disable" -verbose up
	docker-compose -f build/docker-compose.yaml up -d banner_service_app

install:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

dropdb:
	docker exec -it local_pgsql dropdb --username=myuser banner_service

migratedown:
	migrate -path db/migration -database "postgresql://myuser:mypassword@localhost:5432/banner_service?sslmode=disable" -verbose down

