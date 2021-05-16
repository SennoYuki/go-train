// +build wireinject

package main

import (
	"github.com/google/wire"
	"yuki-kabe/internal/data"
)

func initApp() (App, error) {
	panic(wire.Build(data.ProviderSet, newApp))
}
