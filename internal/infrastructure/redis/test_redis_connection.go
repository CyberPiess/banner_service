package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	RedisPassword = "test_password"
)

type TestRedis struct {
	RedisInstance *redis.Client
	container     testcontainers.Container
}

func SetupTestRedis() *TestRedis {

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*60)
	container, redisInstance, err := createContainer(ctx)
	if err != nil {
		log.Fatal("failed to setup test", err)
	}
	cancel()

	return &TestRedis{
		container:     container,
		RedisInstance: redisInstance,
	}
}

func (tr *TestRedis) TearDown() {
	tr.RedisInstance.Close()
	_ = tr.container.Terminate(context.Background())
}

func createContainer(ctx context.Context) (testcontainers.Container, *redis.Client, error) {

	var port = "6379/tcp"

	req := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{port},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}

	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("Could not start redis: %s", err)
	}

	p, err := redisC.MappedPort(ctx, "6379")
	if err != nil {
		log.Fatalf("Could not start redis: %s", err)
	}
	redisAddres := fmt.Sprintf("localhost:%s", p.Port())
	client, err := NewRedis(Config{
		Addres: redisAddres,
	})

	if err != nil {
		return redisC, client, fmt.Errorf("failed to establish redis connection: %v", err)
	}

	return redisC, client, nil
}
