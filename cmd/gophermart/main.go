package main

import (
	"context"

	"github.com/kirilltitov/go-musthave-diploma/internal/app"
	"github.com/kirilltitov/go-musthave-diploma/internal/config"
	"github.com/kirilltitov/go-musthave-diploma/internal/utils"
)

func run() error {
	utils.Log.Infof("Hello")

	cfg := config.New()
	ctx := context.Background()

	a, err := app.New(ctx, cfg, nil)
	if err != nil {
		return err
	}

	utils.Log.Infof("Starting server at %s", cfg.ServerAddress)

	return a.Run()
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}
