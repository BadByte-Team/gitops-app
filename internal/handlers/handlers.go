package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/TU_USUARIO/curso-gitops/internal/auth"
	"github.com/TU_USUARIO/curso-gitops/internal/models"
	"github.com/TU_USUARIO/curso-gitops/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	userRepo     *repository.UserRepository
	moduleRepo   *repository.ModuleRepository
	episodeRepo  *repository.EpisodeRepository
	progressRepo *repository.ProgressRepository
}

func NewHandler(userRepo *repository.UserRepository, moduleRepo *repository.ModuleRepository, episodeRepo *repository.EpisodeRepository, progressRepo *repository.ProgressRepository) *Handler {
	return &Handler{userRepo: userRepo, moduleRepo: moduleRepo, episodeRepo: episodeRepo, progressRepo: progressRepo}
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
	if user.IsBlocked {
		writeJSON(w, http.StatusForbidden, models.APIResponse{Success: false, Message: "tu cuenta ha sido bloqueada"})
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

// POST /api/register
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, models.APIResponse{Success: false, Message: "método no permitido"})
		return
	}
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, models.APIResponse{Success: false, Message: "body inválido"})
		return
	}
	if len(req.Username) < 3 || len(req.Password) < 6 {
		writeJSON(w, http.StatusBadRequest, models.APIResponse{Success: false, Message: "usuario (mín 3 chars) y contraseña (mín 6 chars) requeridos"})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, models.APIResponse{Success: false, Message: "error procesando contraseña"})
		return
	}
	if err := h.userRepo.Create(req.Username, string(hashedPassword), "student"); err != nil {
		writeJSON(w, http.StatusConflict, models.APIResponse{Success: false, Message: "el usuario ya existe o error al crear"})
		return
	}
	writeJSON(w, http.StatusCreated, models.APIResponse{Success: true, Message: "registro exitoso"})
}

// GET /api/modules
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

// POST /api/modules/create
func (h *Handler) CreateModule(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, models.APIResponse{Success: false, Message: "método no permitido"})
		return
	}
	var req models.ModuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, models.APIResponse{Success: false, Message: "body inválido"})
		return
	}
	if req.Title == "" {
		writeJSON(w, http.StatusBadRequest, models.APIResponse{Success: false, Message: "título requerido"})
		return
	}
	id, err := h.moduleRepo.Create(req.Title, req.IsHidden)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, models.APIResponse{Success: false, Message: "error creando módulo"})
		return
	}
	writeJSON(w, http.StatusCreated, models.APIResponse{Success: true, Data: map[string]int64{"id": id}})
}

// PUT /api/modules/update?id=N
func (h *Handler) UpdateModule(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeJSON(w, http.StatusMethodNotAllowed, models.APIResponse{Success: false, Message: "método no permitido"})
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, models.APIResponse{Success: false, Message: "id inválido"})
		return
	}
	var req models.ModuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, models.APIResponse{Success: false, Message: "body inválido"})
		return
	}
	if err := h.moduleRepo.Update(id, req.Title, req.IsHidden); err != nil {
		writeJSON(w, http.StatusInternalServerError, models.APIResponse{Success: false, Message: "error actualizando módulo"})
		return
	}
	writeJSON(w, http.StatusOK, models.APIResponse{Success: true, Message: "módulo actualizado"})
}

// DELETE /api/modules/delete?id=N
func (h *Handler) DeleteModule(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeJSON(w, http.StatusMethodNotAllowed, models.APIResponse{Success: false, Message: "método no permitido"})
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, models.APIResponse{Success: false, Message: "id inválido"})
		return
	}
	if err := h.moduleRepo.Delete(id); err != nil {
		writeJSON(w, http.StatusInternalServerError, models.APIResponse{Success: false, Message: "error eliminando módulo"})
		return
	}
	writeJSON(w, http.StatusOK, models.APIResponse{Success: true, Message: "módulo eliminado"})
}

// GET /api/episodes?module_id=N
func (h *Handler) GetEpisodes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, models.APIResponse{Success: false, Message: "método no permitido"})
		return
	}
	moduleID := r.URL.Query().Get("module_id")
	if moduleID != "" {
		mid, err := strconv.Atoi(moduleID)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, models.APIResponse{Success: false, Message: "module_id inválido"})
			return
		}
		episodes, err := h.episodeRepo.ListByModule(mid)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, models.APIResponse{Success: false, Message: "error consultando episodios"})
			return
		}
		writeJSON(w, http.StatusOK, models.APIResponse{Success: true, Data: episodes})
		return
	}
	episodes, err := h.episodeRepo.ListAll()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, models.APIResponse{Success: false, Message: "error consultando episodios"})
		return
	}
	writeJSON(w, http.StatusOK, models.APIResponse{Success: true, Data: episodes})
}

// POST /api/episodes/create
func (h *Handler) CreateEpisode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, models.APIResponse{Success: false, Message: "método no permitido"})
		return
	}
	var req models.EpisodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, models.APIResponse{Success: false, Message: "body inválido"})
		return
	}
	if req.Title == "" || req.ModuleID <= 0 {
		writeJSON(w, http.StatusBadRequest, models.APIResponse{Success: false, Message: "título y module_id requeridos"})
		return
	}
	id, err := h.episodeRepo.Create(req.ModuleID, req.Title, req.VideoURL, req.IsHidden)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, models.APIResponse{Success: false, Message: "error creando episodio"})
		return
	}
	writeJSON(w, http.StatusCreated, models.APIResponse{Success: true, Data: map[string]int64{"id": id}})
}

// PUT /api/episodes/update?id=N
func (h *Handler) UpdateEpisode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeJSON(w, http.StatusMethodNotAllowed, models.APIResponse{Success: false, Message: "método no permitido"})
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, models.APIResponse{Success: false, Message: "id inválido"})
		return
	}
	var req models.EpisodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, models.APIResponse{Success: false, Message: "body inválido"})
		return
	}
	if err := h.episodeRepo.Update(id, req.Title, req.VideoURL, req.IsHidden); err != nil {
		writeJSON(w, http.StatusInternalServerError, models.APIResponse{Success: false, Message: "error actualizando episodio"})
		return
	}
	writeJSON(w, http.StatusOK, models.APIResponse{Success: true, Message: "episodio actualizado"})
}

// DELETE /api/episodes/delete?id=N
func (h *Handler) DeleteEpisode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeJSON(w, http.StatusMethodNotAllowed, models.APIResponse{Success: false, Message: "método no permitido"})
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, models.APIResponse{Success: false, Message: "id inválido"})
		return
	}
	if err := h.episodeRepo.Delete(id); err != nil {
		writeJSON(w, http.StatusInternalServerError, models.APIResponse{Success: false, Message: "error eliminando episodio"})
		return
	}
	writeJSON(w, http.StatusOK, models.APIResponse{Success: true, Message: "episodio eliminado"})
}

// GET /api/users
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

// PUT /api/users/role?id=N
func (h *Handler) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeJSON(w, http.StatusMethodNotAllowed, models.APIResponse{Success: false, Message: "método no permitido"})
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, models.APIResponse{Success: false, Message: "id inválido"})
		return
	}
	var req models.UserRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, models.APIResponse{Success: false, Message: "body inválido"})
		return
	}
	if req.Role != "admin" && req.Role != "student" {
		writeJSON(w, http.StatusBadRequest, models.APIResponse{Success: false, Message: "rol inválido"})
		return
	}
	if err := h.userRepo.UpdateRole(id, req.Role); err != nil {
		writeJSON(w, http.StatusInternalServerError, models.APIResponse{Success: false, Message: "error actualizando rol"})
		return
	}
	writeJSON(w, http.StatusOK, models.APIResponse{Success: true, Message: "rol actualizado"})
}

// PUT /api/users/block?id=N&blocked=true|false
func (h *Handler) ToggleBlockUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeJSON(w, http.StatusMethodNotAllowed, models.APIResponse{Success: false, Message: "método no permitido"})
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, models.APIResponse{Success: false, Message: "id inválido"})
		return
	}
	blocked := r.URL.Query().Get("blocked") == "true"
	if err := h.userRepo.ToggleBlock(id, blocked); err != nil {
		writeJSON(w, http.StatusInternalServerError, models.APIResponse{Success: false, Message: "error actualizando bloqueo"})
		return
	}
	msg := "usuario desbloqueado"
	if blocked {
		msg = "usuario bloqueado"
	}
	writeJSON(w, http.StatusOK, models.APIResponse{Success: true, Message: msg})
}

// GET /api/progress — autenticado, devuelve episodios completados del usuario
func (h *Handler) GetProgress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, models.APIResponse{Success: false, Message: "método no permitido"})
		return
	}
	userID, _ := strconv.Atoi(r.Header.Get("X-User-ID"))
	completed, err := h.progressRepo.GetCompletedEpisodes(userID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, models.APIResponse{Success: false, Message: "error consultando progreso"})
		return
	}
	if completed == nil {
		completed = []int{}
	}
	writeJSON(w, http.StatusOK, models.APIResponse{Success: true, Data: completed})
}

// POST /api/progress/toggle — autenticado
func (h *Handler) ToggleProgress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, models.APIResponse{Success: false, Message: "método no permitido"})
		return
	}
	userID, _ := strconv.Atoi(r.Header.Get("X-User-ID"))
	var req models.ToggleProgressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, models.APIResponse{Success: false, Message: "body inválido"})
		return
	}
	completed, err := h.progressRepo.ToggleEpisode(userID, req.EpisodeID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, models.APIResponse{Success: false, Message: "error actualizando progreso"})
		return
	}
	writeJSON(w, http.StatusOK, models.APIResponse{Success: true, Data: map[string]bool{"completed": completed}})
}

// GET /api/health
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, models.APIResponse{Success: true, Message: "ok"})
}
