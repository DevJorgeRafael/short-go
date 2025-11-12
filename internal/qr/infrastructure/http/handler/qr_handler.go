package handler

import (
	"fmt"
	"net/http"
	"short-go/config"

	"github.com/go-chi/chi/v5"
	"github.com/skip2/go-qrcode"
)

type QRHandler struct {
	config *config.Config
	// Canal semáforo para controlar el acceso concurrente
	semaphore chan struct{}
}

func NewQRHandler(cfg *config.Config) *QRHandler {
	const maxConcurrentGenerations = 5

	return &QRHandler{
		config:    cfg,
		// Canal con buffer de 5 espacios
		semaphore: make(chan struct{}, maxConcurrentGenerations),
	}
}

func (h *QRHandler) GenerateQR(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	if code == "" {
		http.Error(w, "Código requerido", http.StatusBadRequest)
		return
	}

	baseUrl := h.config.Domain
	if h.config.Port != "" && h.config.Domain == "http://localhost" {
		baseUrl = fmt.Sprintf("%s:%s", h.config.Domain, h.config.Port)
	}
	fullShortUrl := fmt.Sprintf("%s/%s", baseUrl, code)

	// Aplicando concurrencia
	h.semaphore <- struct{}{} // Solicitar acceso

	defer func() { <-h.semaphore }() // Liberar acceso al finalizar

	// Generación del QR
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "public, max-age=86400") // Cache por 1 día

	qr, err := qrcode.New(fullShortUrl, qrcode.Medium)
	if err != nil {
		http.Error(w, "Error generando QR", http.StatusInternalServerError)
		return
	}

	png, err := qr.PNG(256)
	if err != nil {
		http.Error(w, "Error generando QR", http.StatusInternalServerError)
		return
	}

	w.Write(png)
}