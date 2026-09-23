package main

import (
	"net/http"

	"github.com/mucusscraper/gobet-lite/internal/database"
	"github.com/mucusscraper/gobet-lite/internal/server"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		// fornece as dependencias
		fx.Provide(database.ConnectPostgresDB, server.NewServer),
		// aqui, vamos invok
		fx.Invoke(func(mux *http.ServeMux) {}),
	).Run()
}
