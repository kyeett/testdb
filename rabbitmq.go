package testdb

import (
	"fmt"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/streadway/amqp"
)

type RabbitMQContainer struct {
	amqpURL  string
	pool     *dockertest.Pool
	resource *dockertest.Resource
}

func (c *RabbitMQContainer) Connect() (*amqp.Connection, error) {
	var conn *amqp.Connection
	err := c.pool.Retry(func() error {
		var err error
		if conn, err = amqp.Dial(c.amqpURL); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (c *RabbitMQContainer) Close() error {
	return c.pool.Purge(c.resource)
}

func NewRunningRabbitMQContainer() (*RabbitMQContainer, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, err
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "rabbitmq",
		Tag:        "3",
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("amqp://guest:guest@localhost:%s", resource.GetPort("5672/tcp"))
	container := &RabbitMQContainer{
		amqpURL:  url,
		pool:     pool,
		resource: resource,
	}

	return container, nil
}
