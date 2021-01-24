package testdb

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

const (
	password = "123"
	database = "postgres"
)

type PostgresContainer struct {
	databaseURL string
	pool        *dockertest.Pool
	resource    *dockertest.Resource
}

func (c *PostgresContainer) Connect() (*sql.DB, error) {
	var db *sql.DB
	err := c.pool.Retry(func() error {
		var err error
		db, err = sql.Open("postgres", c.databaseURL)
		if err != nil {
			return err
		}
		return db.Ping()
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (c *PostgresContainer) Close() error {
	return c.pool.Purge(c.resource)
}

func (c *PostgresContainer) RunMigrations(migrationsURL string) error {
	m, err := migrate.New(
		migrationsURL,
		c.databaseURL)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		return err
	}

	return nil
}

func NewRunningPostgresContainer() (*PostgresContainer, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, err
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "11.10",
		Env:        []string{"POSTGRES_PASSWORD=" + password, "POSTGRES_DB=" + database},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		return nil, err
	}
	databaseURL := fmt.Sprintf("postgres://postgres:%s@localhost:%s/%s?sslmode=disable", password, resource.GetPort("5432/tcp"), database)

	c := &PostgresContainer{
		databaseURL: databaseURL,
		pool:        pool,
		resource:    resource,
	}
	return c, nil
}
