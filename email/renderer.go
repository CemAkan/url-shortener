package email

import (
	"bytes"
	"github.com/CemAkan/url-shortener/pkg/infrastructure"
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
}

func Render(templateName string, data EmailData) (string, error) {
	files := []string{
		TemplateBasePath + "base.html",
		TemplateBasePath + "components/logo.html",
		TemplateBasePath + "components/header.html", // <== EKLENDİ!
		TemplateBasePath + "components/footer.html",
		TemplateBasePath + "transactional/" + templateName + ".html",
	}

	infrastructure.Log.Infof("Parsing Templates: %s", strings.Join(files, ", "))

	tmpl, err := template.New("base.html").ParseFS(TemplatesFS, files...)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	return buf.String(), err
}
