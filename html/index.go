package html

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"os"
)

//go:embed index.html
var indexHtmlTemplateContent string

func GetIndexHtmlContent(frontendUri string) (string, error) {
	template, err := template.New("").Parse(indexHtmlTemplateContent)
	if err != nil {
		return "", fmt.Errorf("html/GetIndexHtmlContent: could not parse %v template: %w", "index.html", err)
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("html/GetIndexHtmlContent: could not get current working directory: %w", err)
	}

	data := struct {
		HomeDirectory string
		FrontendUri   string
	}{
		HomeDirectory: currentDir,
		FrontendUri:   frontendUri,
	}

	var templateBuffer bytes.Buffer
	err = template.Execute(&templateBuffer, data)
	if err != nil {
		return "", fmt.Errorf("html/GetIndexHtmlContent: could not execute %v template: %w", "index.html", err)
	}

	return templateBuffer.String(), nil
}
