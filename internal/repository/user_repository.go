package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mucusscraper/gobet-lite/internal/models"
)

// aqui vai ter o struct que carrega a conexao com o db
type UserRepository struct {
	db *pgxpool.Pool
}

// aqui uma funcao que vai ter como entrada a conexao e devolve o struct com a propria conexao -
// parece redundante, mas serve para separar os conceitos
func CreateUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// funcao que cria o usuario na tabela e tambem cria o referente na tabela account com o valor 0 no montante
func (r *UserRepository) Create(ctx context.Context, userParams models.CreateUserRequest) (*models.User, error) {
	// query para inserir os valores que nao sao default e retorna o id, usuario e quando foi criado
	queryUser := `
		INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id, username, created_at
	`
	// objeto que vai carregar o id, usuario e quando foi criado na resposta
	res := &models.User{}
	// executa a query e apos isso salva os dados no objeto criado
	err := r.db.QueryRow(ctx, queryUser, userParams.Username, userParams.Password).Scan(&res.ID, &res.Username, &res.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}
	// o mesmo para a tabela account
	queryAccount := `
		INSERT INTO account (user_id) VALUES ($1)
	`
	_, err = r.db.Exec(ctx, queryAccount, res.ID)
	if err != nil {
		return nil, fmt.Errorf("error creating account row for the user: %w", err)
	}
	return res, nil
}
