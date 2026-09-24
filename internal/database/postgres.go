package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

// Lifecycle gerencia como componentes iniciam e terminam gracefully, no caso o componente seria a conexao com o db
func ConnectPostgresDB(lc fx.Lifecycle) (*pgxpool.Pool, error) {
	dbURL := os.Getenv("DATABASE_URL")
	ctx := context.Background()

	// pgxpool é o gerenciador de conexoes da biblioteca pgx, serve para conectar go ao postgres
	// o New cria uma nova conexao
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, fmt.Errorf("error: ", err)
	}

	// a parte do uberfx:
	// o lifecycle é a interface injetada pelo fx
	// o append serve para registrar hooks que o fx vai executar automaticamente quando a aplicacao estiver subindo ou descendo
	// por exemplo, nao precisa escrever manualmente aquele famoso defer db.Close()
	lc.Append(
		fx.Hook{
			// o que deve acontecer ao iniciar o serviço - essa funcao é executada automaticamente logo que voce executa
			// o programa, antes da aplicação começar a receber requisicoes de verdade
			// no caso o ping tenta se comunicar de fato com o banco de dados
			// se falhar, o fx intercepta o erro, trava a subida da aplicação e aborta o programa
			OnStart: func(ctx context.Context) error {
				return pool.Ping(ctx)
			},
			// o que deve acontecer ao desligar o serviço - essa funcao é executada automaticamente ao desligar a aplicacao
			// no caso, vai fechar a conexao com o banco antes de desligar totalmente a aplicacao
			OnStop: func(ctx context.Context) error {
				pool.Close()
				return nil
			},
		})
	return pool, nil
}
