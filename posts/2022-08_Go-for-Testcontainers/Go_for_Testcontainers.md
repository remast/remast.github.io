# Go Integration Tests using Testcontainers

Your application uses a database like PostgreSQL? So how do you test your persistence layer to ensure it's working properly with a real [PostgreSQL](https://www.postgresql.org/) database? Right, you need a test against a real PostgreSQL. Since that test requires external infrastructure it's an _integration test_. You'll learn now how easy it is to write integration tests for external infrastructure using [Testcontainers](https://golang.testcontainers.org/) and [Docker](https://www.docker.com/).

## Integration Test Setup

Our application stores users in a [PostgreSQL](https://www.postgresql.org/) database. It uses a struct `UserRepository` with a method `FindByUsername` that uses plain SQL to find a user by username. We will write an integration test running against a real PostgreSQL in [Docker](https://www.docker.com/) for the method `FindByUsername`.

The integration test for the `FindByUsername` of our `UserRepository` looks like:

```go
func TestUserRepository(t *testing.T) {
	// Setup database
	dbContainer, connPool := SetupTestDatabase()
	defer dbContainer.Terminate(context.Background())

	// Create user repository
	userRepository := NewUserRepository(connPool)

	// Run tests against db
	t.Run("FindExistingUserByUsername", func(t *testing.T) {
		adminUser, err := userRepository.FindByUsername(
			context.Background(),
			"admin",
		)

		is.NoErr(err)
		is.Equal(adminUser.Username, "admin")
	})
}
```

First the database is set up. Then a new `UserRepository` for the test is created with a reference to the connection pool of the database `connPool`. No we run the method to test `userRepository.FindByUsername(ctx, "admin")` and verify the result. But wait, where did that database container come from? Right, we'll set that up using [Testcontainers](https://golang.testcontainers.org/) and [Docker](https://www.docker.com/).

### Database Setup with Testcontainers

We set up the [PostgreSQL](https://www.postgresql.org/) database in a [Docker](https://www.docker.com/) container using the [Testcontainers](https://golang.testcontainers.org/) lib.

As a first step we create a `testcontainers.ContainerRequest` where we set the Docker image to `postgres:14` with exposed port `5432/tcp`. Additionaly the database name as well as username and password are set using environment variables. And to make sure the tests only starts when the database container is up and running we wait for it using the `WaitingFor` option with `wait.ForListeningPort("5432/tcp")`.

Now as step two we run start the requested container.

Finally we use host and port of the running database container to create the connection string for the database with `fmt.Sprintf("postgres://postgres:postgres@%v:%v/testdb", host, port.Port())` and connect with `pgxpool.Connect(context.Background(), dbURI)`.

The whole method `SetupTestDatabase` to set up the PostgreSQL container is (errors omitted):

```go
func SetupTestDatabase() (testcontainers.Container, *pgxpool.Pool) {
	// Request PostgreSQL container
	req := testcontainers.ContainerRequest{
		Image:        "postgres:14",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
		Env: map[string]string{
			"POSTGRES_DB":       "testdb",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_USER":     "postgres",
		},
	}

	// Run PostgreSQL container
	dbContainer, _ := testcontainers.GenericContainer(
		context.Background(),
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
	})

	// Get host and port of PostgreSQL container
	host, _ := dbContainer.Host(context.Background())
	port, _ := dbContainer.MappedPort(context.Background(), "5432")

	// Create db connection string and connect
	dbURI := fmt.Sprintf("postgres://postgres:postgres@%v:%v/testdb", host, port.Port())
	connPool, _ := pgxpool.Connect(context.Background(), dbURI)

	return dbContainer, connPool
}
```

Notice that we make sure the PostgreSQL container is terminated after our integration tests with `defer dbContainer.Terminate(context.Background())`.

### Migration

### Wrap Upa
