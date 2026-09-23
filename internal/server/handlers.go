package server

import (
	"encoding/json"
	"net/http"

	"github.com/mucusscraper/gobet-lite/internal/models"
	"github.com/mucusscraper/gobet-lite/internal/repository"
)

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
	_, err = h.UserRepo.Create(r.Context(), *userRequest)
	if err != nil {
		http.Error(w, "erro interno ao criar usuario (pode ter outro usuario com mesmo nome)", http.StatusInternalServerError)
		return
	}

	// resposta http
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
