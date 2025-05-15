package email

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
)

type EmailData struct {
	Title            string
	VerificationLink string
}

func Render(templateName string, data EmailData) (string, error) {
	files := []string{
		TemplateBasePath + "base.html",
		TemplateBasePath + "components/logo.html",
		TemplateBasePath + "components/footer.html",
		TemplateBasePath + "transactional/" + templateName + ".html",
	}

	// Debug log:
	fmt.Println("🟢 Parsing Templates:", strings.Join(files, ", "))

	tmpl, err := template.New("base.html").ParseFS(TemplatesFS, files...)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	return buf.String(), err
}
