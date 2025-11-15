package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"short-go/internal/auth/application/service"
	"time"
)

type BrevoEmailService struct {
	apiKey     string
	senderEmail string
	httpClient *http.Client
}

var _ service.EmailService = (*BrevoEmailService)(nil)

func NewBrevoEmailService(apiKey string, senderEmail string) *BrevoEmailService {
	return &BrevoEmailService{
		apiKey: apiKey,
		senderEmail: senderEmail,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

const brevoApiUrl = "https://api.brevo.com/v3/smtp/email"

type brevoEmailSender struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
type brevoEmailRecipient struct {
	Email string `json:"email"`
}
type brevoSendEmailPayload struct {
	Sender      brevoEmailSender      `json:"sender"`
	To          []brevoEmailRecipient `json:"to"`
	Subject     string                `json:"subject"`
	HtmlContent string                `json:"htmlContent"`
}

func (s *BrevoEmailService) SendPasswordResetCode(toEmail string, code string) error {
	payload := brevoSendEmailPayload{
		Sender: brevoEmailSender{
			Name:  "ShortGo Support",
			Email: s.senderEmail,
		}, 
		To: []brevoEmailRecipient{
			{Email: toEmail},
		},
		Subject:     "Tu código para recuperar la contraseña",
		HtmlContent: fmt.Sprintf(
			"<h1>Recuperación de Contraseña</h1><p>Tu código de un solo uso es:</p><h2>%s</h2><p>Este código expira en 10 minutos.</p>",
			code,
		),
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error al 'marshal' el payload de Brevo: %w", err)
	}

	req, err := http.NewRequest("POST", brevoApiUrl, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("error al crear la solicitud HTTP a Brevo: %w", err)
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", s.apiKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error al enviar la solicitud HTTP a Brevo: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("la API de Brevo devolvió un error (stattus %d)", resp.StatusCode)
	}

	return nil
}