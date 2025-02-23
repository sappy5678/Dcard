package redis_test

import (
	"context"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/sappy5678/dcard/pkg/utl/redis"
)

func TestNew(t *testing.T) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "redis/redis-stack:7.4.0-v3",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForAll(wait.ForLog("Ready to accept connections"), wait.ForListeningPort("6379")),
	}
	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(t, err)

	endpoint, err := redisC.Endpoint(ctx, "")
	assert.NoError(t, err)

	defer testcontainers.CleanupContainer(t, redisC)

	_, err = redis.New("mockAddress")
	if err == nil {
		t.Error("Expected error")
	}

	client, err := redis.New(endpoint)
	if err != nil {
		t.Fatalf("Error establishing connection %v", err)
	}

	assert.NotNil(t, client)
}
