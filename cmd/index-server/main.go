package main

import (
	"github.com/fmoura/index-api/internal/handlers"
	"gofr.dev/pkg/gofr"
)

func main() {
	// initialise gofr object
	app := gofr.New()

	// register route index
	app.GET("/index/{index}", handlers.HandleIndex)

	// Runs the server, it will listen on the default port 8000.
	// it can be over-ridden through configs
	app.Run()
}
