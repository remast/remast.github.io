---
title: "Go for the REST"
date: 2019-09-01T23:11:13Z
draft: false
---

The Go programming language is a great fit for RESTful web applications. Picking the right Go web framework to start with is not an easy task. Lucky enough it is a task solved already. Pick the Go web framework [buffalo](https://gobuffalo.io) and off you go. Not convinced yet? Learn how to build a RESTful web application with [buffalo](https://gobuffalo.io).

[Buffalo](https://gobuffalo.io) is a whole ecosystem to develop web applications in Go. Buffalo combines routing, templating and testing to build web applications based on the famous [model-view-controller pattern](https://en.wikipedia.org/wiki/Model%E2%80%93view%E2%80%93controller). The foundation is a powerful tooling with hot reload, a generator for views and controller and tools for common tasks like compressing JavaScript or database migration. Buffalo combines existing libraries like the Gorilla web toolkit with best practices of Go web development. With buffalo comes a uniform and consistent project structure covering packages and directories up until handling the database session. So let's get started.

## Installing Buffalo
We'll use the new [Go modules](https://blog.golang.org/using-go-modules) for dependency management. So let's activate that using the environment variable `GO111MODULE=on`. Now we can install buffalo with `go get -u -v github.com/gobuffalo/buffalo/buffalo`. Executing the command `buffalo` in your shell should display the help.

## Create and run the RESTful web application
Now we create our first buffalo RESTful web application named `rowdy`. The command `buffalo new rowdy --api` creates a new RESTful buffalo web application (without frontend) in the directory `rowdy`. We can run the application with the command `buffalo dev` (from directory `rowdy`). The application is available at http://localhost:3000. But wait, that doesn't work! The generated application needs a running database. so we start a PostgreSQL database with `docker run --name rowdy_db -e POSTGRES_PASSWORD=postgres -d postgres`. Now you should be able to start the application and access the REST API under the web root saying `"Welcome to Buffalo"`.

You can keep the application running while you develop. Buffalo will automatically pick up changes in sources and rebuild and restart the application as needed.

## Guide through the App
Entry point of the buffalo application is the `main.go` in the root directory of our application. The main.go is the main package of the application and starts the application by executing the method `Serve` of the buffalo struct [`buffalo.App`](https://godoc.org/github.com/gobuffalo/buffalo#App). The buffalo app itself is created in the `app.go` which is in the package actions.

The app.go is where the music plays. In the `app.go` all routes and middlewares are registered. For every route an action is registered. An action is the controller layer of the MVC pattern. For the moment we only have one action which is the `HomeHandler`. With `app.GET("/", HomeHandler)` the action HomeHandler is registered in the app.go for the path `/`. The `HomeHandler` itself is declared in home.go in the package actions, and the associated test is in `home_test.go`.

The package model is for the model layer of MVC. Buffalo assumes the model layer is the persistence layer of the application running against a relational database like postgres. This model layer is based on  buffalo pop, a active record inspired persistence layer for Go.

The directory migrations is all about maintaining the schema of the relational database. Here scripts for the database migration framework [fizz](https://gobuffalo.io/en/docs/db/fizz) are stored. But we'll get to that later.

Then there are the directories fixtures, config and grifts. Fixtures are scripts to create database content for tests. Config is for the configuration of the buffalo code generator. Grifts are small scripts to automate common tasks like listing all routes.

To complete our guide we look at the generated configuration files. The most important config file of our application is the database.yml with the database credentials for several environments. Then there is the Dockerfile to create a docker container of our app.

## Create a Health Check Action
Now let's create our first action, a health check at route `/health/check`. We generate the skeleton of our health check action using the buffalo code generator by running `buffalo generate action health check --skip-template`. That generated an action with name `HealthCheck`. The first parameter of the generator command is the name of the action and the second parameter check is the name of the handler. Since we build a RESTful API we don't need a HTML template for our action and suppress generating one with the parameter `--skip-template`. Buffalo created the files `health.go` and `health_test.go` in the directory actions. The generator also registered our new action for the path `/health/check` in the `app.go`.

A buffalo action needs to implement the interface [`buffalo.Handler`](https://godoc.org/github.com/gobuffalo/buffalo#Handler). The terms action and handler are both synonymously used for the controller layer of the MVC pattern. A look at the `home.go` shows the generated code:

```go
func HealthCheck(c buffalo.Context) error {
    return c.Render(200, r.HTML("health/check.html"))
}
```

The generated handler tries to render the template `health/check.html` as HTML. Since we build a RESTful API we don't have templates so this can't work. We change the implementation to `c.Render(200, r.String("Up and running!"))` so that it returns a plain text string. 

```go
func HealthCheck(c buffalo.Context) error {
	return c.Render(200, r.String("Up and running!"))
}
```

## Test the Health Check
Now we need to fix the test for our health check. The whole buffalo application is available it the test. In the test you can run requests against the application or even access the database it's running against of the HTTP session.

To test the health check we run request path `/health/check` and verify the HTTP status code. The code to do that is:

```go
package actions

func (as *ActionSuite) Test_Health_Check() {
	as.Fail("Not Implemented!")
	result := as.HTML("/health").Get()

	as.Equal(200, result.Code)
	as.Contains(result.Body.String(), "Up and running")
}
```

## More about Actions and Routes
When configuring routes you can use placeholders and regular expressions in the path like `app.GET("/account/{category}/{id:[0-9]+}", AccountHandler)`. Since that's base on the Gorilla web toolkit you can make full use of the Gorilla multiplexer.

In the handler itself you have full access to the [`buffalo.Context`](https://godoc.org/github.com/gobuffalo/buffalo#Context). With the buffalo context you can render JSON or HTML and access path and query parameters with them methods `c.Params()` and `c.Param(string)`. You can also acess the logger and log health check requests with `c.Logger().Info("Health checked by: %v", c.Request().Host)`.

To get an overview of all registered routes use the command `buffalo routes` which prints:

```bash
METHOD | PATH           | ALIASES | NAME            | HANDLER
------ | ----           | ------- | ----            | -------
GET    | /              |         | rootPath        | rowdy/actions.HomeHandler
GET    | /health/check/ |         | healthCheckPath | rowdy/actions.HealthCheck
```

### CRUD Resources
To make it easy to develop CRUD functionality buffalo has resources. A buffalo resource always contains the actions Create, Show, Update, Destroy and List with standardized paths in the router. Resources are always RESTful and never use templates.

Now we create a resource to manage concerts with `buffalo generate resource concert`. Buffalo created several files for us. The implemntation of the concert resource in the files `concert.go` in the package actions, along with the tests in `concert_test.go`. A database mapping for pop in `concert.go` in the package models including test. And last but not least in the directory migrations two fizz scripts to create a table for the concerts in our database schema. The concerts resources is registered in the `app.go` using `app.Resource("/concerts", ConcertsResource{})`.

Looking at the output of buffalo routes we see all the routes of our new concert resource.

```bash
METHOD | PATH                    | ALIASES | NAME            | HANDLER
------ | ----                    | ------- | ----            | -------
GET    | /                       |         | rootPath        | rowdy/actions.HomeHandler
GET    | /concerts/              |         | concertsPath    | rowdy/actions.ConcertsResource.List
POST   | /concerts/              |         | concertsPath    | rowdy/actions.ConcertsResource.Create
GET    | /concerts/{concert_id}/ |         | concertPath     | rowdy/actions.ConcertsResource.Show
PUT    | /concerts/{concert_id}/ |         | concertPath     | rowdy/actions.ConcertsResource.Update
DELETE | /concerts/{concert_id}/ |         | concertPath     | rowdy/actions.ConcertsResource.Destroy
GET    | /health/check/          |         | healthCheckPath | rowdy/actions.HealthCheck
```

Let's now look into the generated handler method `Show`. We see that buffalo gets the database connection from the buffalo context with `c.Value("tx").(*pop.Connection)`. Then buffalo loads the concert model from the database using pop. The id of the concert model is provided as a path parameter. Finally the loaded concert model is returned with `c.Render(200, r.Auto(c, concert))`. Depending on the negotiated content-type this is either as JSON or as XML.

```go
// GET /concerts/{concert_id}
func (v ConcertsResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Concert
	concert := &models.Concert{}

	// To find the Concert the parameter concert_id is used.
	if err := tx.Find(concert, c.Param("concert_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.Auto(c, concert))
}
```

## Error Handling
Handling errors properly can be tricky. But buffalo has you covered. Errors generally occur in handlers. Handlers have two options of dealing with them:

1. Return a normal Go error from the handler like `errors.New("Boom!")`. The application then returns HTTP status code 500 for internal server error.

2. Use the bufallo context method Error to handle like `return context.Error(501, errors.New("That blew up!"))`. That way the applications returns the custom status code 501.

Either way the application keeps running and the error is logged and returned either as JSON or XML. If buffalo runs in dev mode (environment `variable GO_ENV="development"` or start with `buffalo dev`) this will contain a stacktrace of the error.

## Middlewares
Buffalo applications feature many best practices in Go web development. This is mainly featured by various middleware. A middleware is code that runs befor or after a HTTP handler (in buffalo called action). Middlewares are a general concept of Go web development. A good read about the concept of middleware is [Making and Using HTTP Middleware](https://www.alexedwards.net/blog/making-and-using-middleware) by Alex Edwards.

Buffalo comes with lots of usefule middlewares integrated into the [`buffalo.App`](https://godoc.org/github.com/gobuffalo/buffalo#App) right from the start. E.g. middleware that prevents your app from crashing after a panic or an error. In our generated app middleware is registered in `actions/app.go`. One example is `app.Use(forceSSL)` which registers a middleware that redirects HTTP requests to HTTPS.

It's also possible to skip middleware for specific routes. We could use this to skip authentication middleware for our health check route. The command `buffalo task middleware` shows the registered middlewares for every route.

## Configuration
A buffalo app is configured using environment variables as proposed in [the twelve-factor app](https://12factor.net/config). The configuration is powered by the buffalop package [envy](https://github.com/gobuffalo/envy).

There are predefined environment variables like `GO_ENV` for the current environment (e.g. development, test, production) or `HOST` (`localhost`) and `PORT` (3000). You can also use custom variables. You can use an optional environment variable with `envy.Get("CUSTOM_SETTING", "default-value")`. Or a mandatory variable with `envy.MustGet("MY_REQUIRED_OPTION")` which results in an error if the variable is not set.

You can also read environment variables from the file `.env`. Buffalo looks for a file named `.env` in the current directory and evaluates it. This nicely integrates with [direnv](https://direnv.net/) if you create a `.envrc` with only the line `dotenv` (see [`.env` files](https://github.com/direnv/direnv/issues/284)). You can use this to have buffalo read your local development settings which of course you should never check in.

## Plugins
Buffalo can be extended using plugins. The database framework pop is implemented as a buffalo plugin. There are also plugins for [kubernetes](https://toolkit.gobuffalo.io/tools/69219770-87c4-4e53-9b6e-92f81ab819e4/) or [Azure](https://toolkit.gobuffalo.io/tools/d299aa2d-03be-4df6-bd05-9b19d3f19ae6/). The buffalo website contains a [list of available plugins](https://toolkit.gobuffalo.io/tools?topic=plugin).

## Build the App
To build the app run `buffalo build`. Buffalo compiles the app and creates a standalone binary `rowdy` in the directory `build`.

## Wrap up the REST
We learned how to build a RESTful application with the buffalo framework that adheres to the MVC pattern. The model layer consists of structs mapped against a PostgreSQL database using pop. Actions and resources make up the controller layer. We created a simple health check action and a CRUD resource for concerts. We have not used the view layer yet since all we need is serializing structs to XML or JSON. We learned how to handle error, how to deal with configuration and how to run a production build in the CI environment.

Buffalo comes with a uniform and consistent structure for web applications in go, follows best practices of Go web development and covers the full lifecycle from development to production.

Yet there's still so much that we could't cover. Buffalo has a lot to offer to build frontends. Or you can use buffalos tasks or background workers. All this is covered by the great [buffalo documentation](https://gobuffalo.io/en/docs/index/).

And of course stay tuned because I'll spread more buffalo love here soon!