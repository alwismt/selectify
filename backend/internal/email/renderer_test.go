package email

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRender_LoginNotification(t *testing.T) {
	html, err := Render(TemplateLoginNotification, struct {
		FirstName   string
		LoggedInAt  string
		IP          string
		UserAgent   string
		Country     string
		City        string
		Subdivision string
		Timezone    string
	}{
		FirstName:   "Ada",
		LoggedInAt:  "2026-08-06 19:00:00 UTC",
		IP:          "203.0.113.10",
		UserAgent:   "Mozilla/5.0",
		Country:     "United States",
		City:        "New York",
		Subdivision: "New York",
		Timezone:    "America/New_York",
	})
	require.NoError(t, err)
	require.Contains(t, html, "Selectify")
	require.Contains(t, html, "Hi Ada,")
	require.Contains(t, html, "203.0.113.10")
	require.Contains(t, html, "Mozilla/5.0")
	require.Contains(t, html, "2026-08-06 19:00:00 UTC")
	require.Contains(t, html, "United States")
	require.Contains(t, html, "America/New_York")
	require.True(t, strings.Contains(html, "<!DOCTYPE html>") || strings.Contains(html, "<html"))
}

func TestRender_MissingTemplate(t *testing.T) {
	_, err := Render("does_not_exist", nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "does_not_exist")
}

func TestRender_EmptyName(t *testing.T) {
	_, err := Render("", nil)
	require.Error(t, err)
}
