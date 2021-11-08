package main

import (
	"net/http"

	"github.com/SedaOzy/go-getir-case-study/configuration"
	"github.com/SedaOzy/go-getir-case-study/handlers"
)

type App struct {
	Config *configuration.Config
	Router *http.ServeMux
}

func (a *App) Init() {
	cfg, isSuccessful := configuration.Init()
	if !isSuccessful {
		return
	}

	a.Config = cfg
	a.Router = handlers.InitRouter(a.Config)
}

func (a *App) Run() {
	handlers.Run(a.Config, a.Router)
}

func main() {
	a := App{}
	a.Init()
	a.Run()
}
