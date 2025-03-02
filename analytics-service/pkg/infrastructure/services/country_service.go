package services

import (
	"github.com/oschwald/geoip2-golang"
	"go.uber.org/zap"
	"net"
)

type CountryService struct {
	logger *zap.Logger
}

func NewCountryService(logger *zap.Logger) *CountryService {
	return &CountryService{
		logger: logger,
	}
}

func (s *CountryService) GetCountryByIP(ip string) string {
	db, err := geoip2.Open("./data/GeoLite2-Country.mmdb")
	if err != nil {
		s.logger.Error("Error opening GeoIP database", zap.Error(err))
		return "Unknown"
	}
	defer db.Close()

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return "Unknown"
	}

	record, err := db.Country(parsedIP)
	if err != nil {
		return "Unknown"
	}

	countryName := record.Country.Names["en"]
	if countryName == "" {
		return "Unknown"
	}

	return countryName
}
