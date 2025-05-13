package utils

import "fmt"

func GenerateResetPasswordEmail(username, resetLink string) string {
	return fmt.Sprintf(`
		<h2>Password Reset Request</h2>
		<p>Hello %s,</p>
		<p>Click the link below to reset your password:</p>
		<a href="%s">%s</a>
		<p>This link will expire in 15 minutes.</p>
	`, username, resetLink, resetLink)
}

func GenerateEmailVerification(username, verifyLink string) string {
	return fmt.Sprintf(`
		<h2>Email Verification</h2>
		<p>Welcome %s!</p>
		<p>Please verify your email by clicking the link below:</p>
		<a href="%s">%s</a>
	`, username, verifyLink, verifyLink)
}
