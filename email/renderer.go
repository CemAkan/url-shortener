package email

import (
	"bytes"
	"github.com/CemAkan/url-shortener/pkg/infrastructure/logger"
	"html/template"
	"strings"
)

type EmailData struct {
	Title            string
	Greeting         string
	Message          string
	VerificationLink string
	LogoURL          string
	HeaderURL        string
	ButtonText       string
}

func Render(data EmailData) (string, error) {
	files := []string{
		TemplateBasePath + "base.html",
		TemplateBasePath + "components/logo.html",
		TemplateBasePath + "components/header.html",
		TemplateBasePath + "components/footer.html",
		TemplateBasePath + "transactional/content.html",
	}

	logger.Log.Infof("Parsing Templates: %s", strings.Join(files, ", "))

	tmpl, err := template.New("base.html").ParseFS(TemplatesFS, files...)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	return buf.String(), err
}
