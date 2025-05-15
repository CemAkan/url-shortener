package email

import (
	"bytes"
	"html/template"
)

type EmailData struct {
	Title            string
	VerificationLink string
}

func Render(templateName string, data EmailData) (string, error) {
	tmpl, err := template.New("base.html").ParseFS(TemplatesFS,
		"templates/base.html",
		"templates/components/logo.html",
		"templates/components/footer.html",
		"templates/transactional/"+templateName+".html",
	)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	return buf.String(), err
}
