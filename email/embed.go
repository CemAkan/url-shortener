package email

import (
	"embed"
	"io/fs"
)

//go:embed templates/**/*.html
var embeddedTemplates embed.FS

// TemplatesFS exposes embedded templates with FS interface.
var TemplatesFS fs.FS = embeddedTemplates
