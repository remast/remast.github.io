---
title: Go Integration Tests using Testcontainers
published: true
description: How do you test your Go application is working properly with a real PostgreSQL database? Right, you need to test against a real db. You'll learn how easy it is to write integration tests for your Go application using Testcontainers and Docker.
tags: go, beginners
date: 2022-08-29T12:11:13Z
draft: false
canonicalUrl: "https://dev.to/remast/go-integration-tests-using-testcontainers-9o5"
# Use a ratio of 100:42 for best results.
---

Your application uses a database like PostgreSQL? So how do you test your persistence layer to ensure it's working properly with a real [PostgreSQL](https://www.postgresql.org/) database? Right, you need to test against a real PostgreSQL. Since that test requires external infrastructure it's an _integration test_. You'll learn now how easy it is to write integration tests for external infrastructure using [Testcontainers](https://golang.testcontainers.org/) and [Docker](https://www.docker.com/).

## Integration Test Setup

Our application stores users in a [PostgreSQL](https://www.postgresql.org/) database. It uses the struct `UserRepository` with a method `FindByUsername` that uses plain SQL to find a user by username. We will write an integration test running against a real PostgreSQL in [Docker](https://www.docker.com/) for the method `FindByUsername`.

The integration test for the `FindByUsername` of our `UserRepository` looks like:

```go
func TestUserRepository(t *testing.T) {
	// Setup database
	dbContainer, connPool := SetupTestDatabase()
	defer dbContainer.Terminate(context.Background())

	// Create user repository
	repository := NewUserRepository(connPool)

	// Run tests against db
	t.Run("FindExistingUserByUsername", func(t *testing.T) {
		adminUser, err := repository.FindByUsername(
			context.Background(),
			"admin",
		)

		is.NoErr(err)
		is.Equal(adminUser.Username, "admin")
	})
}
```

First the database is set up. Then a new `UserRepository` is created for the test with a reference to the connection pool of the database `connPool`. No we run the method to test `userRepository.FindByUsername(ctx, "admin")` and verify the result. But wait, where did that database container come from? Right, we'll set that up using [Testcontainers](https://golang.testcontainers.org/) and [Docker](https://www.docker.com/).

## Database Setup using Testcontainers

We set up the [PostgreSQL](https://www.postgresql.org/) database in a [Docker](https://www.docker.com/) container using the [Testcontainers](https://golang.testcontainers.org/) lib.

As a first step we create a `testcontainers.ContainerRequest` where we set the Docker image to `postgres:14` with exposed port `5432/tcp`. Additionaly the database name as well as username and password are set using environment variables. And to make sure the tests only starts when the database container is up and running we wait for it using the `WaitingFor` option with `wait.ForListeningPort("5432/tcp")`.

Now as second step we start the requested container.

Finally in step 3 we use host and port of the running database container to create the connection string for the database with `fmt.Sprintf("postgres://postgres:postgres@%v:%v/testdb", host, port.Port())`. Now we connect with `pgxpool.Connect(context.Background(), dbURI)`.

The whole method `SetupTestDatabase` to set up the PostgreSQL container is (errors omitted):

```go
func SetupTestDatabase() (testcontainers.Container, *pgxpool.Pool) {
	// 1. Create PostgreSQL container request
	containerReq := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
		Env: map[string]string{
			"POSTGRES_DB":       "testdb",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_USER":     "postgres",
		},
	}

	// 2. Start PostgreSQL container
	dbContainer, _ := testcontainers.GenericContainer(
		context.Background(),
		testcontainers.GenericContainerRequest{
			ContainerRequest: containerReq,
			Started:          true,
	})

	// 3.1 Get host and port of PostgreSQL container
	host, _ := dbContainer.Host(context.Background())
	port, _ := dbContainer.MappedPort(context.Background(), "5432")

	// 3.2 Create db connection string and connect
	dbURI := fmt.Sprintf("postgres://postgres:postgres@%v:%v/testdb", host, port.Port())
	connPool, _ := pgxpool.Connect(context.Background(), dbURI)

	return dbContainer, connPool
}
```

Notice that we make sure the PostgreSQL container is terminated after our integration tests with `defer dbContainer.Terminate(context.Background())`.

## Adding Database Migrations

So far our test starts out with an empty database. That's not very useful since we need the database tables of our application. In our simple example we need the table `users`. We will now set up our database using [golang-migrate](https://github.com/golang-migrate/migrate).

We add the database migrations to the `SetupTestDatabase()` method by adding the call `MigrateDb(dbURI)`.

```go
func SetupTestDatabase() (testcontainers.Container, *pgxpool.Pool) {
	// ... (see before)

	// 3.2 Create db connection string and connect
	dbURI := fmt.Sprintf("postgres://postgres:postgres@%v:%v/testdb", host, port.Port())
	MigrateDb(dbURI)
	connPool, _ := pgxpool.Connect(context.Background(), dbURI)

	return dbContainer, connPool
}
```

The method `MigrateDb(dbURI)` applies the database migration scripts to the database using [golang-migrate](https://github.com/golang-migrate/migrate). The migration scripts are read from the directory `migrations` which is embedded into the binary of our application.

```go
//go:embed migrations
var migrations embed.FS

func MigrateDb(dbURI string) error {
	source, _ := iofs.New(migrations, "migrations")

	m, _ := migrate.NewWithSourceInstance("iofs", source, strings.Replace(dbURI, "postgres://", "pgx://", 1))
	defer m.Close()

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
```

## Wrap Up

We have a working setup for integration tests against a real [PostgreSQL](https://www.postgresql.org/) database running in [Docker](https://www.docker.com/) using [Testcontainers](https://golang.testcontainers.org/). We can use this setup for integration tests of our persistence layer. But that's not all it's good for.

This setup is a great way for all kinds of integration tests that need infrastructure. E.g. a test that send emails to an mail server running in docker, as in [mail_resource_smtp_test.go](https://github.com/Baralga/baralga-app/blob/main/shared/mail_resource_smtp_test.go).

Code see [github.com/remast/go-for-testcontainers](https://github.com/remast/go-for-testcontainers)
