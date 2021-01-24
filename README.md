# Example usage

## With `github.com/stretchr/testify`
#### Setup test suite
```go
type ExampleTestSuite struct {
	suite.Suite
	db          *sql.DB
	pgContainer *testdb.PostgresContainer
}

func (s *ExampleTestSuite) SetupSuite() {
	// Create & run
	pgContainer, err := testdb.NewRunningPostgresContainer()
	require.NoError(s.T(), err)

	// Connect
	db, err := pgContainer.Connect()
	require.NoError(s.T(), err)

	// Migrate database
	err = pgContainer.RunMigrations("file://db/migrations")
	require.NoError(s.T(), err)

	s.db = db
	s.pgContainer = pgContainer
}

func (s *ExampleTestSuite) TearDownSuite() {
	// Make sure to clear everything after test suite
	assert.NoError(s.T(), s.db.Close())
	assert.NoError(s.T(), s.pgContainer.Close(), "could not purge resource")
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, &ExampleTestSuite{})
}
```

Postgres connection can be used in tests like this:
```go
// Test pinging our db
func (s *ExampleTestSuite) TestPingDB() {
	err := s.db.Ping()
	s.NoError(err, "failed to ping db")
}
```