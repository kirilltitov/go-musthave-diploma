package gophermart

import (
	"github.com/kirilltitov/go-musthave-diploma/internal/config"
	"github.com/kirilltitov/go-musthave-diploma/internal/container"
)

type Gophermart struct {
	config    config.Config
	container *container.Container
}

func New(cfg config.Config, cnt *container.Container) Gophermart {
	return Gophermart{config: cfg, container: cnt}
}
