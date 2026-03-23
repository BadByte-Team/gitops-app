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
	// Reintentar conexión a la BD (MySQL puede tardar en arrancar)
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
	h := handlers.NewHandler(userRepo, moduleRepo)

	mux := http.NewServeMux()

	// Rutas públicas
	mux.HandleFunc("/api/health", h.Health)
	mux.HandleFunc("/api/login", h.Login)

	// Rutas autenticadas
	mux.HandleFunc("/api/modules", auth.RequireAuth(h.GetModules))

	// Rutas solo admin
	mux.HandleFunc("/api/users", auth.RequireAdmin(h.GetUsers))

	// Servir archivos estáticos del frontend
	fs := http.FileServer(http.Dir("./frontend"))
	mux.Handle("/", fs)

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
