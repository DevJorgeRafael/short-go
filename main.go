package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"short-go/config"
	"short-go/internal/shared/infrastructure"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error cargando config:", err)
	}

	db, err := config.InitDatabase(cfg.DatabaseUrl)
	if err != nil {
		log.Fatal("Error inicializando base de datos:", err)
	}
	log.Println("Base de datos conectada")

	// Dependency Injection Container
	container := infrastructure.NewContainer(db, cfg)

	r := chi.NewRouter()

	// Configuración de CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"}, 
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300, // Cachear la respuesta (preflight) por 5 min
	}))
	
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

	log.Printf("Servidor escuchando en %s\n", cfg.Domain+addr)
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
