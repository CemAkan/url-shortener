package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"time"
)

type EmailData struct {
	Title       string
	Greeting    string
	Message     string
	ButtonText  string
	ButtonLink  string
	ShowButton  bool
	CompanyName string
	Year        int
}

// RenderEmailTemplate renders the HTML email with given data.
func RenderEmailTemplate(data EmailData) (string, error) {
	const tpl = `
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>{{.Title}}</title>
	<style>
		body { font-family: Arial, sans-serif; background: #f9f9f9; color: #333; padding: 20px; }
		.container { max-width: 600px; margin: auto; background: #fff; border-radius: 8px; box-shadow: 0 0 10px rgba(0,0,0,0.1); overflow: hidden; }
		.header { background: #007bff; color: white; text-align: center; padding: 20px; }
		.content { padding: 20px; }
		.button { display: inline-block; padding: 12px 24px; background: #28a745; color: white; text-decoration: none; border-radius: 4px; margin-top: 20px; }
		.footer { background: #f0f0f0; text-align: center; font-size: 12px; color: #888; padding: 10px; }
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1>{{.Title}}</h1>
		</div>
		<div class="content">
			<p>{{.Greeting}}</p>
			<p>{{.Message}}</p>
			{{ if .ShowButton }}
			<a href="{{.ButtonLink}}" class="button">{{.ButtonText}}</a>
			{{ end }}
			<p style="font-size: 12px; color: #999; margin-top: 30px;">If you did not request this, please ignore this email.</p>
		</div>
		<div class="footer">
			&copy; {{.Year}} {{.CompanyName}}. All rights reserved.
		</div>
	</div>
</body>
</html>
`

	tmpl, err := template.New("email").Parse(tpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// GenerateResetPasswordEmail returns password reset email body
func GenerateResetPasswordEmail(username, resetLink string) string {
	data := EmailData{
		Title:       "Password Reset Request",
		Greeting:    fmt.Sprintf("Hello %s,", username),
		Message:     "You requested to reset your password. Click the button below. This link will expire in 15 minutes.",
		ButtonText:  "Reset Password",
		ButtonLink:  resetLink,
		ShowButton:  true,
		CompanyName: "CemAkan",
		Year:        time.Now().Year(),
	}

	body, _ := RenderEmailTemplate(data)
	return body
}

// GenerateEmailVerification returns email verification body
func GenerateEmailVerification(username, verifyLink string) string {
	data := EmailData{
		Title:       "Email Verification",
		Greeting:    fmt.Sprintf("Welcome %s!", username),
		Message:     "Please verify your email by clicking the button below.",
		ButtonText:  "Verify Email",
		ButtonLink:  verifyLink,
		ShowButton:  true,
		CompanyName: "CemAkan",
		Year:        time.Now().Year(),
	}

	body, _ := RenderEmailTemplate(data)
	return body
}
