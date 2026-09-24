package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mucusscraper/gobet-lite/internal/websockethub"
	"go.uber.org/fx"
)

type ServerParams struct {
	// esse fx.In vai servir para que quando essa funcao for citada no main, nao precisar passar um objeto com essa struct
	// como entrada na funcao, o proprio fx vai olhar para dentro do container e ver o que pede.
	// ou seja, caça as dependencias no container e preenche a struct automaticamente antes de chamar a funcao
	// no caso o lifecycle é nativo do fx e o *pgxpool.Pool vai vir do ConnectPostgresDB
	fx.In
	Lifecycle   fx.Lifecycle
	DB          *pgxpool.Pool
	UserHandler *UserHandler
	BetHandler  *BetHandler
	WsHub       *websockethub.Hub
}

// funcao simples de health
func HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func NewServer(p ServerParams) *http.ServeMux {
	// inicia o server mux que vai servir para registrar as rotas
	mux := http.NewServeMux()

	// vamos criar uma rota de health simples
	mux.HandleFunc("/health", HandleHealth)
	mux.HandleFunc("POST /users", p.UserHandler.CreateUser)
	mux.HandleFunc("POST /bets", p.BetHandler.CreateBet)
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) { ServeWs(p.WsHub, w, r) })

	// cria os parametros para o servidor - a porta e o server mux com as rotas registradas
	// porem agora quem vai iniciar mesmo o servidor vai ser o UberFX
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	// nova parte do lifecycle
	p.Lifecycle.Append(
		// hook para iniciar e desligar
		fx.Hook{
			// no start aqui vamos deixar o listenandserve rodando e se falhar -> precisa de go func pq é bloqueante
			OnStart: func(ctx context.Context) error {
				go func() {
					fmt.Println("servidor rodando na porta :8080...")
					if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
						fmt.Printf("erro no servidor HTTP:%v\n", err)
					}
				}()
				return nil
			},
			// no stop vamos desligar o servidor
			OnStop: func(ctx context.Context) error {
				fmt.Println("desligando o servidor...")
				return server.Shutdown(ctx)
			},
		})
	return mux
}
