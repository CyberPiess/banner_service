include ./build/.env

all:
	docker-compose -f build/docker-compose.yaml up -d cache
	docker-compose -f build/docker-compose.yaml up -d db
	docker exec -it local_pgsql createdb --username=${POSTGRES_USER} --owner=${POSTGRES_USER} ${DBNAME}
	migrate -path db/migration -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:5432/${DBNAME}?sslmode=${SSLMODE}" -verbose up
	docker-compose -f build/docker-compose.yaml up -d banner_service_app

install:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

dropdb:
	docker exec -it local_pgsql dropdb --username=${POSTGRES_USER} ${DBNAME}

migratedown:
	migrate -path db/migration -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:5432/${DBNAME}?sslmode=${SSLMODE}" -verbose down

