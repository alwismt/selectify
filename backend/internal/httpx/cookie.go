package httpx

import (
	"net/http"
	"time"

	"alwis.dev/selectify/internal/model"
)

const (
	Cookie1      = "slf"     // session bearer secret
	Cookie2      = "slf_ss"  // signed-in flag
	DeviceCookie = "slf_did" // non-authenticating device identifier
)

func SetSessionCookies(s *model.UserSession, w http.ResponseWriter) {
	maxAge := 0
	if s.RememberMe {
		maxAge = int(time.Until(s.ExpiresAt).Seconds())
		if maxAge < 1 {
			maxAge = int((30 * 24 * time.Hour).Seconds())
		}
	}

	http.SetCookie(w, createCookie(Cookie1, s.RawSessionToken, true, true, true, maxAge))
	http.SetCookie(w, createCookie(Cookie2, "true", true, true, true, maxAge))

	if s.RawDeviceToken != "" {
		deviceMaxAge := int((365 * 24 * time.Hour).Seconds())
		http.SetCookie(w, createCookie(DeviceCookie, s.RawDeviceToken, true, true, true, deviceMaxAge))
	}
}

func DeleteSessionCookies(w http.ResponseWriter) {
	http.SetCookie(w, deleteCookie(Cookie1, "", true, true, true))
	http.SetCookie(w, createCookie(Cookie2, "false", true, true, true, 0))
}

func createCookie(n string, v string, secure, sameSite, httpOnly bool, maxAge int) *http.Cookie {
	c := &http.Cookie{
		Name:     n,
		Value:    v,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: httpOnly,
		Secure:   secure,
	}

	if sameSite {
		c.SameSite = http.SameSiteStrictMode
	} else {
		c.SameSite = http.SameSiteNoneMode
	}

	return c
}

func deleteCookie(n string, v string, secure, sameSite, httpOnly bool) *http.Cookie {
	c := &http.Cookie{
		Name:     n,
		Value:    v,
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: httpOnly,
		Secure:   secure,
	}

	if sameSite {
		c.SameSite = http.SameSiteStrictMode
	} else {
		c.SameSite = http.SameSiteNoneMode
	}
	return c
}
