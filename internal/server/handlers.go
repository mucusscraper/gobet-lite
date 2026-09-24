package server

import (
	"encoding/json"
	"net/http"

	"github.com/mucusscraper/gobet-lite/internal/models"
	"github.com/mucusscraper/gobet-lite/internal/repository"
)

// USER HANDLER

type UserHandler struct {
	UserRepo *repository.UserRepository
}

func NewUserHandler(UserRepo *repository.UserRepository) *UserHandler {
	return &UserHandler{
		UserRepo: UserRepo,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// verifica se o metodo é post
	if r.Method != http.MethodPost {
		http.Error(w, "Metodo nao permitido", http.StatusMethodNotAllowed)
		return
	}

	// cria o objeto de DTO contendo username e password
	userRequest := &models.CreateUserRequest{}

	// decodifica o payload que chegou no objeto de DTO
	err := json.NewDecoder(r.Body).Decode(userRequest)
	if err != nil {
		http.Error(w, "Payload json invalido", http.StatusBadRequest)
		return
	}

	// uma validacao simples
	if userRequest.Username == "" || userRequest.Password == "" {
		http.Error(w, "Username e password sao obrigatorios", http.StatusBadRequest)
		return
	}

	// cria na tabela o usuario
	user, err := h.UserRepo.Create(r.Context(), *userRequest)
	if err != nil {
		http.Error(w, "erro interno ao criar usuario (pode ter outro usuario com mesmo nome)", http.StatusInternalServerError)
		return
	}
	// resposta http
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(user)
}

// BET HANDLER

type BetHandler struct {
	BetRepo *repository.BetRepository
}

func NewBetHandler(BetRepo *repository.BetRepository) *BetHandler {
	return &BetHandler{
		BetRepo: BetRepo,
	}
}

func (b *BetHandler) CreateBet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Metodo nao permitido", http.StatusMethodNotAllowed)
		return
	}
	betRequest := &models.CreateBetRequest{}
	err := json.NewDecoder(r.Body).Decode(betRequest)
	if err != nil || betRequest.Input <= 0 || betRequest.Multiplier <= 0 {
		http.Error(w, "Dados da aposta invalidos", http.StatusBadRequest)
		return
	}
	bet, err := b.BetRepo.CreateBet(r.Context(), *betRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(bet)
}
