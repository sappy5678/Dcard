package postgres_test

import (
	"testing"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"

	"github.com/sappy5678/dcard/pkg/utl/postgres"
)

func TestNew(t *testing.T) {
	pgdb := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Password("password").
		Port(3001).Logger(nil))

	err := pgdb.Start()
	assert.NoError(t, err)
	defer pgdb.Stop()

	_, err = postgres.New("PSN")
	if err == nil {
		t.Error("Expected error")
	}

	_, err = postgres.New("postgres://postgres:password@localhost:1234/postgres?sslmode=disable")
	if err == nil {
		t.Error("Expected error")
	}

	db, err := postgres.New("postgres://postgres:password@localhost:3001/postgres?sslmode=disable")
	if err != nil {
		t.Fatalf("Error establishing connection %v", err)
	}

	db.Ping()

	assert.NotNil(t, db)

	db.Close()
}
