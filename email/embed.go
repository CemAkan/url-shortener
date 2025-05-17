package email

import (
	"embed"
	"github.com/CemAkan/url-shortener/pkg/logger"
)

//go:embed **/base.html
var baseFile embed.FS

var TemplateBasePath string

func init() {
	candidates := []string{
		"templates/base.html",
		"email/templates/base.html",
	}

	for _, path := range candidates {
		_, err := baseFile.Open(path)
		if err == nil {
			TemplateBasePath = path[:len(path)-len("base.html")]
			break
		}
	}

	if TemplateBasePath == "" {
		logger.Log.Fatalf("Could not detect template base path")
	}
}

//go:embed **/*
var TemplatesFS embed.FS
