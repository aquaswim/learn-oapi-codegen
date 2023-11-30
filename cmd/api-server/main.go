package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	middleware "github.com/oapi-codegen/echo-middleware"
	"net"
	"os"
	"os/signal"
	apiServer "todo-codegen/internal/api_server"
)

func main() {
	port := flag.String("port", "8080", "Port for test HTTP server")
	flag.Parse()

	swagger, err := apiServer.GetSwagger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading swagger spec\n: %s", err)
		os.Exit(1)
	}

	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swagger.Servers = nil

	// Create an instance of our handler which satisfies the generated interface
	apiHandlers := apiServer.New()

	// This is how you set up a basic Echo router
	e := echo.New()
	// Log all requests
	e.Use(echomiddleware.Logger())
	// Recovery middleware
	e.Use(echomiddleware.Recover())
	// Error handler
	e.Use(apiServer.ErrorHandler)
	// Use our validation middleware to check all requests against the
	// OpenAPI schema.
	e.Use(middleware.OapiRequestValidator(swagger))

	// We now register our petStore above as the handler for the interface
	apiServer.RegisterHandlers(e, apiHandlers)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		e.Logger.Infof("Gracefully shutting down...")
		err := e.Shutdown(context.Background())
		if err != nil {
			e.Logger.Errorf("error shutdown server %s", err)
		}
	}()

	// And we serve HTTP until the world ends.
	_ = e.Start(net.JoinHostPort("0.0.0.0", *port))
	// do other cleanup here
}
