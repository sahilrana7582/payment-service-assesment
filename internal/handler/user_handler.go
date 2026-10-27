package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/sahilrana27582/payment-system/internal/service"
)

type UserCreateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user UserCreateRequest

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user.Username = strings.TrimSpace(user.Username)
	user.Email = strings.TrimSpace(user.Email)

	if user.Username == "" || user.Password == "" || user.Email == "" {
		http.Error(w, "Username, password, and email are required", http.StatusBadRequest)
		return
	}

	userID, err := h.svc.CreateUser(r.Context(), user.Username, user.Password, user.Email)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEmailTaken):
			http.Error(w, "Email is already registered", http.StatusConflict)
		case errors.Is(err, service.ErrPasswordTooLong):
			http.Error(w, "Password must be at most 72 bytes", http.StatusBadRequest)
		default:
			log.Printf("failed to create user: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(map[string]string{
		"id": userID,
	}); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}
