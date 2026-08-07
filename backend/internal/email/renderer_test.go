package email

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const testLogoURL = "https://cdn.example/logos/logo.svg"

func TestRender_LoginNotification(t *testing.T) {
	html, err := Render(TemplateLoginNotification, struct {
		FirstName  string
		LastName   string
		LoggedInAt string
		IP         string
		UserAgent  string
		Location   string
	}{
		FirstName:  "Ada",
		LastName:   "Lovelace",
		LoggedInAt: "2026-08-06 19:00:00 UTC",
		IP:         "203.0.113.10",
		UserAgent:  "Mozilla/5.0",
		Location:   "New York, United States",
	}, testLogoURL)
	require.NoError(t, err)
	require.Contains(t, html, `src="`+testLogoURL+`"`)
	require.Contains(t, html, "Selectify")
	require.Contains(t, html, "Hi Ada Lovelace,")
	require.Contains(t, html, "203.0.113.10")
	require.Contains(t, html, "Mozilla/5.0")
	require.Contains(t, html, "2026-08-06 19:00:00 UTC")
	require.Contains(t, html, "New York, United States")
	require.True(t, strings.Contains(html, "<!DOCTYPE html>") || strings.Contains(html, "<html"))
}

func TestRender_PasswordReset(t *testing.T) {
	html, err := Render(TemplatePasswordReset, struct {
		FirstName   string
		Location    string
		IP          string
		RequestedAt string
		ResetURL    string
	}{
		FirstName:   "Ada",
		Location:    "Kaunas, Lithuania",
		IP:          "84.xxx.xxx.xxx",
		RequestedAt: "August 7, 2026, 01:15",
		ResetURL:    "http://localhost:3000/reset-password?token=abc",
	}, testLogoURL)
	require.NoError(t, err)
	require.Contains(t, html, `src="`+testLogoURL+`"`)
	require.Contains(t, html, "We received a request to reset your password.")
	require.Contains(t, html, "This link expires in 5 minutes.")
	require.Contains(t, html, "Kaunas, Lithuania")
	require.Contains(t, html, "84.xxx.xxx.xxx")
	require.Contains(t, html, "reset-password?token=abc")
}

func TestRender_PasswordChanged(t *testing.T) {
	html, err := Render(TemplatePasswordChanged, struct {
		FirstName string
		Location  string
		IP        string
		ChangedAt string
	}{
		FirstName: "Ada",
		Location:  "Kaunas, Lithuania",
		IP:        "84.xxx.xxx.xxx",
		ChangedAt: "August 7, 2026, 01:15",
	}, testLogoURL)
	require.NoError(t, err)
	require.Contains(t, html, `src="`+testLogoURL+`"`)
	require.Contains(t, html, "password was changed successfully")
	require.Contains(t, html, "Kaunas, Lithuania")
}

func TestRender_MissingTemplate(t *testing.T) {
	_, err := Render("does_not_exist", nil, testLogoURL)
	require.Error(t, err)
	require.Contains(t, err.Error(), "does_not_exist")
}

func TestRender_EmptyName(t *testing.T) {
	_, err := Render("", nil, testLogoURL)
	require.Error(t, err)
}
