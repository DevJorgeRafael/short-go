package main

import (
	"context"
	"fmt"
	"go-task-easy-list/config"
	"go-task-easy-list/internal/shared/infrastructure"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	// "github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error cargando config:", err)
	}

	db, err := config.InitDatabase(cfg.DBPath)
	if err != nil {
		log.Fatal("Error inicializando base de datos:", err)
	}
	log.Println("Base de datos conectada")

	// Dependency Injection Container
	container := infrastructure.NewContainer(db, cfg.JWTSecret)

	r := chi.NewRouter()
	// r.Use(middleware.Logger)
	// r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Registrar todas las rutas de los módulos
	container.RegisterRoutes(r)

	// Server
	addr := fmt.Sprintf(":%s", cfg.Port)
	server := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go gracefulShutdown(server)

	log.Printf("Servidor escuchando en http://localhost%s\n", addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Error del servidor:", err)
	}
}

func gracefulShutdown(server *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	
	<-quit
	log.Println("Apagando servidor...")
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Error al apagar servidor:", err)
	}
	
	log.Println("Servidor detenido correctamente")
}