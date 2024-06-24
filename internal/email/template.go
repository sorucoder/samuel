package email

import (
	"embed"
	"html/template"
	"path/filepath"

	"github.com/Masterminds/sprig/v3"
)

//go:embed templates/*.go.html
var templatesFS embed.FS

var (
	PasswordChangeRequestTemplate *Template = newTemplate("Password Change Request", "password_change_request.go.html")
)

type Template struct {
	subject string
	raw     *template.Template
}

func newTemplate(subject string, name string) *Template {
	return &Template{
		subject: subject,
		raw: template.Must(
			template.
				New(name).
				Funcs(sprig.FuncMap()).
				ParseFS(templatesFS, filepath.Join("templates", name), "templates/common.go.html"),
		),
	}
}
