package main

import (
	"net/http"

	"github.com/mucusscraper/gobet-lite/internal/database"
	"github.com/mucusscraper/gobet-lite/internal/repository"
	"github.com/mucusscraper/gobet-lite/internal/server"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		// fornece as dependencias
		fx.Provide(database.ConnectPostgresDB, repository.CreateUserRepository, server.NewUserHandler, repository.NewBetRepository, server.NewBetHandler, server.NewServer),
		// aqui, vamos invok
		fx.Invoke(func(mux *http.ServeMux) {}),
	).Run()
}
