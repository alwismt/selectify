package service

import (
	"context"
	"net"

	"github.com/kelseyhightower/envconfig"
	"github.com/oschwald/geoip2-golang"

	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/program"
)

type Location struct {
	Country     string
	City        string
	Subdivision string
	Timezone    string
}

type GeoIPService interface {
	Lookup(ip string) Location
	Close() error
}

type geoIPConfig struct {
	DBPath string `envconfig:"db_path" required:"true"`
}

type geoIPService struct {
	db *geoip2.Reader
}

func NewGeoIPService() GeoIPService {
	var wrapper struct {
		GeoIP geoIPConfig `envconfig:"geoip"`
	}

	prefix := program.AppPrefix
	if prefix == "" {
		prefix = "EVT"
	}

	if err := envconfig.Process(prefix, &wrapper); err != nil {
		panic(logger.Error(context.Background(), err, "failed to process GeoIP env vars"))
	}

	db, err := geoip2.Open(wrapper.GeoIP.DBPath)
	if err != nil {
		panic(logger.Error(context.Background(), err, "failed to open GeoIP database"))
	}

	return &geoIPService{db: db}
}

func (s *geoIPService) Lookup(ipStr string) Location {
	if s == nil || s.db == nil || ipStr == "" {
		return Location{}
	}

	ip := net.ParseIP(ipStr)
	if ip == nil || !isPublicIP(ip) {
		return Location{}
	}

	record, err := s.db.City(ip)
	if err != nil || record == nil {
		return Location{}
	}

	loc := Location{
		Country:  englishName(record.Country.Names, record.Country.IsoCode),
		City:     englishName(record.City.Names, ""),
		Timezone: record.Location.TimeZone,
	}
	if len(record.Subdivisions) > 0 {
		sub := record.Subdivisions[0]
		loc.Subdivision = englishName(sub.Names, sub.IsoCode)
	}
	return loc
}

func (s *geoIPService) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

func isPublicIP(ip net.IP) bool {
	return !ip.IsLoopback() &&
		!ip.IsPrivate() &&
		!ip.IsUnspecified() &&
		!ip.IsLinkLocalUnicast() &&
		!ip.IsLinkLocalMulticast() &&
		!ip.IsMulticast()
}

func englishName(names map[string]string, fallback string) string {
	if name := names["en"]; name != "" {
		return name
	}
	return fallback
}
