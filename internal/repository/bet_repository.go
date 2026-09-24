package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mucusscraper/gobet-lite/internal/models"
	"github.com/mucusscraper/gobet-lite/internal/websockethub"
)

type BetRepository struct {
	db    *pgxpool.Pool
	wsHub *websockethub.Hub
}

func NewBetRepository(db *pgxpool.Pool, ws *websockethub.Hub) *BetRepository {
	return &BetRepository{
		db:    db,
		wsHub: ws,
	}
}

func (r *BetRepository) CreateBet(ctx context.Context, betParams models.CreateBetRequest) (*models.Bet, error) {
	// aqui vamos fazer no modelo de transaction para garantir a atomicidade das bets e impedir que, se algo falhar,
	// o banco cancele tudo e ninguem saia perdendo nada
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro ao iniciar transação das bets: %w", err)
	}
	// garante que se houver erro, a transação será revertida
	defer tx.Rollback(ctx)

	// verificar o saldo atual e se a conta está bloqueada

	var currentBalance float64
	var isBlocked bool
	fmt.Printf("%v", betParams.UserID)
	err = tx.QueryRow(ctx, `SELECT total, blocked FROM account WHERE user_id=$1 FOR UPDATE`, betParams.UserID).Scan(&currentBalance, &isBlocked)
	fmt.Printf("%v", currentBalance)
	fmt.Printf("%v", isBlocked)
	fmt.Printf("%v", betParams.UserID)

	if err != nil {
		// Retorne o erro real para debuggar mais fácil
		return nil, fmt.Errorf("erro ao buscar conta do usuário: %w", err)
	}
	if isBlocked {
		return nil, errors.New("conta bloqueada")
	}
	if currentBalance < betParams.Input {
		return nil, errors.New("saldo insuficiente")
	}

	result := betParams.Input * betParams.Multiplier
	newBalance := currentBalance - betParams.Input + result

	// inserir os novos dados na tabela account
	_, err = tx.Exec(ctx, `UPDATE account SET total=$1 WHERE user_id=$2`, newBalance, betParams.UserID)
	if err != nil {
		return nil, errors.New("erro ao atualizar saldo")
	}

	// inserir nova row em bets
	var bet models.Bet
	queryBet := `
    INSERT INTO bets (user_id, input, multiplier, result) 
    VALUES ($1, $2, $3, $4) 
    RETURNING id, user_id, input, multiplier, result, created_at
`
	err = tx.QueryRow(ctx, queryBet, betParams.UserID, betParams.Input, betParams.Multiplier, result).Scan(&bet.ID, &bet.UserID, &bet.Input, &bet.Multiplier, &bet.Result, &bet.CreatedAt)
	if err != nil {
		return nil, errors.New("erro ao registrar aposta")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro ao commitar transação: %w", err)
	}
	if r.wsHub != nil {
		r.wsHub.SendToUser(bet.UserID, map[string]interface{}{
			"event":       "bet.result",
			"bet_id":      bet.ID,
			"input":       bet.Input,
			"multiplier":  bet.Multiplier,
			"result":      bet.Result,
			"new_balance": newBalance,
		})
	}
	return &bet, nil
}
