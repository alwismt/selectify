package service

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/oschwald/geoip2-golang"
	"github.com/stretchr/testify/require"
)

func openTestGeoIP(t *testing.T) *geoIPService {
	t.Helper()
	path := filepath.Join("..", "data", "GeoLite2-City.mmdb")
	if _, err := os.Stat(path); err != nil {
		t.Skipf("GeoLite2 DB not available at %s: %v", path, err)
	}
	db, err := geoip2.Open(path)
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	return &geoIPService{db: db}
}

func TestGeoIPService_Lookup_PrivateAndInvalid(t *testing.T) {
	svc := openTestGeoIP(t)

	require.Equal(t, Location{}, svc.Lookup(""))
	require.Equal(t, Location{}, svc.Lookup("not-an-ip"))
	require.Equal(t, Location{}, svc.Lookup("127.0.0.1"))
	require.Equal(t, Location{}, svc.Lookup("10.0.0.1"))
	require.Equal(t, Location{}, svc.Lookup("192.168.1.1"))
	require.Equal(t, Location{}, svc.Lookup("::1"))
}

func TestGeoIPService_Lookup_PublicIP(t *testing.T) {
	svc := openTestGeoIP(t)

	loc := svc.Lookup("8.8.8.8")
	require.NotEmpty(t, loc.Country)
	require.NotEmpty(t, loc.Timezone)
}
