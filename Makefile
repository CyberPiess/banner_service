cache:
	docker-compose -f build/docker-compose.yaml up -d cache

install:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

postgres:
	docker-compose -f build/docker-compose.yaml up -d db

createdb:
	docker exec -it local_pgsql createdb --username=myuser --owner=myuser banner_service

dropdb:
	docker exec -it local_pgsql dropdb --username=myuser banner_service

migrateup:
	migrate -path db/migration -database "postgresql://myuser:mypassword@localhost:5432/banner_service?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://myuser:mypassword@localhost:5432/banner_service?sslmode=disable" -verbose down

app:
	docker-compose -f build/docker-compose.yaml up -d banner_service_app
