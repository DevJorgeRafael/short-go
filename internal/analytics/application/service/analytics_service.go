package service

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	analyticsModel "short-go/internal/analytics/domain/model"
	analyticsRepo "short-go/internal/analytics/domain/repository"
	shortLinkRepo "short-go/internal/short-links/domain/repository"
	"strings"
	"time"
)

var (
	ErrUnauthorized = errors.New("no autorizado para ver estas analíticas")
	ErrLinkNotFound = errors.New("enlace no encontrado")
)

type AnalyticsService struct {
	clickRepo   analyticsRepo.ClickRepository
	shortLinkRepo shortLinkRepo.ShortLinkRepository
	clickChannel chan *analyticsModel.Click
}

func NewAnalyticsService(clickRepo analyticsRepo.ClickRepository, shortLinkRepo shortLinkRepo.ShortLinkRepository) *AnalyticsService {
	s := &AnalyticsService{
		clickRepo:     clickRepo,
		shortLinkRepo: shortLinkRepo,
		// Buffer de 100 clicks para aguantar picos de tráfico
		clickChannel: make(chan *analyticsModel.Click, 100),
	}

	// Worker en segundo plano
	go s.processClicks()

	return s
}

//  --------------- FUNCIONALIDAD DE REGISTRAR  ---------------
func (s *AnalyticsService) TrackClick(code, ip, userAgen, referrer string) {
	// Aqui se llamará a una API de GeoIP para obtener el countryCode real
	countryCode := "XX"

	click := &analyticsModel.Click{
		LinkCode:    code,
		IPAddress: ip,
		UserAgent:  userAgen,
		Referrer:    referrer,
		CountryCode: countryCode,
		ClickedAt: time.Now(),
	}

	// Enviar al canal para procesamiento asíncrono
	select {
	case s.clickChannel <- click:
		// Click enviado al canal
	default:
		log.Println("Warning: Analytics buffer full, dropping click")
	}
}

func (s *AnalyticsService) processClicks() {
	for click := range s.clickChannel {
		if click.CountryCode == "XX" {
			click.CountryCode = s.resolveCountryCode(click.IPAddress)
		}

		if err := s.clickRepo.Save(click); err != nil {
			log.Printf("Error saving click: %v", err)
		}
	}
}

func (s *AnalyticsService) GetStats(code string, managementToken string, userID *string) (*analyticsModel.LinkStats, error) {
	link, err := s.shortLinkRepo.FindByCode(code)
	if err != nil {
		return nil, ErrLinkNotFound
	}

	isAuthorized := false

	if userID != nil && link.UserID != nil && *userID == *link.UserID {
		isAuthorized = true
	}

	if !isAuthorized && link.ManagementToken == managementToken {
		isAuthorized = true
	}

	if !isAuthorized {
		return nil, ErrUnauthorized
	}

	return s.clickRepo.GetLinkStats(code)
}


// ---------------- FUNCIONES AUXILIARES  ---------------
func (s *AnalyticsService) resolveCountryCode(ip string) string {
	if strings.Contains(ip, "127.0.0.1") || strings.Contains(ip, "::1") {
		return "EC"
	}

	client := http.Client{Timeout: 2 * time.Second}
	url := "http://ip-api.com/json/" + ip + "?fields=countryCode"
	
	resp, err := client.Get(url)
	if err != nil {
		return "XX" //Error de red
	}
	defer resp.Body.Close()

	var result struct {
		CountryCode string `json:"countryCode"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "XX" //Error de parseo
	}

	if result.CountryCode == "" {
		return "XX" //No se encontró el país
	}

	return result.CountryCode
}