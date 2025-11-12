package config

import (
	"short-go/config"
	"short-go/internal/qr/infrastructure/http/handler"

	"github.com/go-chi/chi/v5"
)

type QRModule struct {
	Handler *handler.QRHandler
}

func NewQRModule(cfg *config.Config) *QRModule {
	qrHandler := handler.NewQRHandler(cfg)

	return &QRModule{
		Handler: qrHandler,
	}
}

func (m *QRModule) RegisterRoutes(r chi.Router) {
	r.Get("/api/qr/{code}", m.Handler.GenerateQR)
}