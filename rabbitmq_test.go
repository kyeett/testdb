package testdb

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRabbitMQ(t *testing.T) {
	rabbitContainer, err := NewRunningRabbitMQContainer()
	require.NoError(t, err)

	// Connect
	conn, err := rabbitContainer.Connect()
	require.NoError(t, err)

	assert.NoError(t, conn.Close())
	assert.NoError(t, rabbitContainer.Close(), "could not purge resource")
}

