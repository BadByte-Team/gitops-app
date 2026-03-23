package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/TU_USUARIO/curso-gitops/internal/auth"
	"github.com/TU_USUARIO/curso-gitops/internal/models"
	"github.com/TU_USUARIO/curso-gitops/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	userRepo   *repository.UserRepository
	moduleRepo *repository.ModuleRepository
}

func NewHandler(userRepo *repository.UserRepository, moduleRepo *repository.ModuleRepository) *Handler {
	return &Handler{userRepo: userRepo, moduleRepo: moduleRepo}
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// POST /api/login
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, models.APIResponse{Success: false, Message: "método no permitido"})
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, models.APIResponse{Success: false, Message: "body inválido"})
		return
	}

	user, err := h.userRepo.FindByUsername(req.Username)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, models.APIResponse{Success: false, Message: "credenciales inválidas"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		writeJSON(w, http.StatusUnauthorized, models.APIResponse{Success: false, Message: "credenciales inválidas"})
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, models.APIResponse{Success: false, Message: "error generando token"})
		return
	}

	writeJSON(w, http.StatusOK, models.APIResponse{
		Success: true,
		Data:    models.LoginResponse{Token: token, Role: user.Role},
	})
}

// GET /api/modules — autenticado
func (h *Handler) GetModules(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, models.APIResponse{Success: false, Message: "método no permitido"})
		return
	}

	role := r.Header.Get("X-Role")

	var modules interface{}
	var err error

	if role == "admin" {
		modules, err = h.moduleRepo.ListAll()
	} else {
		modules, err = h.moduleRepo.ListVisible()
	}

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, models.APIResponse{Success: false, Message: "error consultando módulos"})
		return
	}

	writeJSON(w, http.StatusOK, models.APIResponse{Success: true, Data: modules})
}

// GET /api/users — solo admin
func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, models.APIResponse{Success: false, Message: "método no permitido"})
		return
	}

	users, err := h.userRepo.ListAll()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, models.APIResponse{Success: false, Message: "error consultando usuarios"})
		return
	}

	writeJSON(w, http.StatusOK, models.APIResponse{Success: true, Data: users})
}

// GET /api/health — sin auth, para readiness check
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, models.APIResponse{Success: true, Message: "ok"})
}
