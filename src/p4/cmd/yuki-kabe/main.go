package main

import (
	"yuki-kabe/internal/data"
)

func main() {
	initApp()
}

type App struct {
	d *data.Data
}

func newApp(d *data.Data) App {
	return App{d: d}
}
