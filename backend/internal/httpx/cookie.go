package httpx

import (
	"net/http"

	"alwis.dev/selectify/internal/model"
)

const (
	Cookie1 = "slf"
	Cookie2 = "slf_ss"
)

func SetSessionCookies(s *model.UserSession, w http.ResponseWriter) {
	http.SetCookie(w, createCookie(Cookie1, s.SessionId.String(), true, true, true))
	http.SetCookie(w, createCookie(Cookie2, "true", true, true, true))
}

func DeleteSessionCookies(w http.ResponseWriter) {
	http.SetCookie(w, deleteCookie(Cookie1, "", true, true, true))
	http.SetCookie(w, createCookie(Cookie2, "false", true, true, true))
}

func createCookie(n string, v string, secure, sameSite, httpOnly bool) *http.Cookie {
	c := &http.Cookie{
		Name:     n,
		Value:    v,
		Path:     "/",
		MaxAge:   0,
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
