package email

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
)

//go:embed templates/*.html
var templateFS embed.FS

const (
	TemplateLoginNotification = "login_notification"
	TemplatePasswordReset     = "password_reset"
	TemplatePasswordChanged   = "password_changed"
)

var templates *template.Template

func init() {
	templates = template.Must(template.ParseFS(templateFS, "templates/*.html"))
}

// Render executes a named body template and wraps it in the shared layout.
func Render(name string, data any) (string, error) {
	if name == "" {
		return "", fmt.Errorf("template name is required")
	}
	if name == "layout" {
		return "", fmt.Errorf("template %q is not a content template", name)
	}

	var bodyBuf bytes.Buffer
	if err := templates.ExecuteTemplate(&bodyBuf, name+".html", data); err != nil {
		return "", fmt.Errorf("render template %q: %w", name, err)
	}

	var out bytes.Buffer
	if err := templates.ExecuteTemplate(&out, "layout.html", map[string]any{
		"Body": template.HTML(bodyBuf.String()),
	}); err != nil {
		return "", fmt.Errorf("render layout: %w", err)
	}

	return out.String(), nil
}
