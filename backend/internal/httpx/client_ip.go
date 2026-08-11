package httpx

import (
	"net"
	"net/http"
	"strings"
)

func ClientIP(r *http.Request) string {
	remoteIP := r.RemoteAddr

	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		remoteIP = host
	}

	// Use with Cloudflare
	if cfIP := net.ParseIP(strings.TrimSpace(r.Header.Get("CF-Connecting-IP"))); cfIP != nil {
		return cfIP.String()
	}

	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		return xff
	}

	if xrip := net.ParseIP(strings.TrimSpace(r.Header.Get("X-Real-IP"))); xrip != nil {
		return xrip.String()
	}

	return remoteIP
}
