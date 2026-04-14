package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/TU_USUARIO/curso-gitops/internal/auth"
	"github.com/TU_USUARIO/curso-gitops/internal/handlers"
	"github.com/TU_USUARIO/curso-gitops/internal/repository"
)

func main() {
	var db *repository.DB
	var err error
	for i := 0; i < 10; i++ {
		db, err = repository.NewDB()
		if err == nil {
			break
		}
		log.Printf("Intento %d: error conectando a BD: %v — reintentando en 3s...", i+1, err)
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		log.Fatalf("No se pudo conectar a la base de datos: %v", err)
	}
	defer db.Close()
	log.Println("Conexión a BD establecida")

	userRepo := repository.NewUserRepository(db)
	moduleRepo := repository.NewModuleRepository(db)
	episodeRepo := repository.NewEpisodeRepository(db)
	progressRepo := repository.NewProgressRepository(db)
	h := handlers.NewHandler(userRepo, moduleRepo, episodeRepo, progressRepo)

	mux := http.NewServeMux()

	// Rutas públicas
	mux.HandleFunc("/api/health", h.Health)
	mux.HandleFunc("/api/login", h.Login)
	mux.HandleFunc("/api/register", h.Register)

	// Rutas autenticadas
	mux.HandleFunc("/api/modules", auth.RequireAuth(h.GetModules))
	mux.HandleFunc("/api/episodes", auth.RequireAuth(h.GetEpisodes))
	mux.HandleFunc("/api/progress", auth.RequireAuth(h.GetProgress))
	mux.HandleFunc("/api/progress/toggle", auth.RequireAuth(h.ToggleProgress))

	// Rutas admin — Módulos CRUD
	mux.HandleFunc("/api/modules/create", auth.RequireAdmin(h.CreateModule))
	mux.HandleFunc("/api/modules/update", auth.RequireAdmin(h.UpdateModule))
	mux.HandleFunc("/api/modules/delete", auth.RequireAdmin(h.DeleteModule))

	// Rutas admin — Episodios CRUD
	mux.HandleFunc("/api/episodes/create", auth.RequireAdmin(h.CreateEpisode))
	mux.HandleFunc("/api/episodes/update", auth.RequireAdmin(h.UpdateEpisode))
	mux.HandleFunc("/api/episodes/delete", auth.RequireAdmin(h.DeleteEpisode))

	// Rutas admin — Gestión de usuarios
	mux.HandleFunc("/api/users", auth.RequireAdmin(h.GetUsers))
	mux.HandleFunc("/api/users/role", auth.RequireAdmin(h.UpdateUserRole))
	mux.HandleFunc("/api/users/block", auth.RequireAdmin(h.ToggleBlockUser))

	// Frontend (URLs limpias)
	servePage := func(file string) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "./frontend/"+file)
		}
	}
	mux.HandleFunc("/", servePage("index.html"))
	mux.HandleFunc("/register", servePage("register.html"))
	mux.HandleFunc("/dashboard", servePage("dashboard.html"))
	mux.HandleFunc("/admin", servePage("admin.html"))

	port := getEnv("PORT", "8080")
	addr := fmt.Sprintf(":%s", port)

	log.Printf("Servidor GitOps iniciado en %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Error iniciando servidor: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
