package main

import (
	"fmt"
	"strconv"

	"github.com/fmoura/index-server/internal/data"
	"github.com/fmoura/index-server/internal/handler"
	"github.com/fmoura/index-server/internal/service"
	"gofr.dev/pkg/gofr"
)

const (
	conformationLevelKey     = "CONFORMATION_LEVEL"
	defaultConformationValue = "10"
)

func main() {
	// initialise gofr object
	app := gofr.New()

	// Create Data Service
	dataProvider, err := data.NewTextDataProvider(app.Logger())
	if err != nil {
		panic(err)
	}

	conformValue, err := strconv.ParseUint(app.Config.GetOrDefault(conformationLevelKey, defaultConformationValue), 10, 8)
	if err != nil {
		panic(fmt.Errorf("Error reading %s : %w", conformationLevelKey, err))
	}
	if conformValue > 100 {
		panic(fmt.Errorf("Invalid value for %s. It must be 0-100", conformationLevelKey))
	}
	dataService := service.NewIndexService(app.Logger(), dataProvider, conformValue)

	// register route index
	indexHandler := handlers.NewIndexHandler(dataService)
	app.GET(handlers.IndexValuePath, indexHandler.HandleGet)

	// Runs the server, it will listen on the default port 8000.
	// it can be over-ridden through configs
	app.Run()
}
