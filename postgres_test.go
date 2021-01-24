package testdb

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPostgres(t *testing.T) {
	pgContainer, err := NewRunningPostgresContainer()
	require.NoError(t, err)

	// Connect
	db, err := pgContainer.Connect()
	require.NoError(t, err)

	assert.NoError(t, db.Close())
	assert.NoError(t, pgContainer.Close(), "could not purge resource")
}
