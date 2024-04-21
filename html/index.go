package html

import (
	"bytes"
	"context"
	_ "embed"
	"ffmpeg-clipper/html/templ"
	"fmt"
	"html/template"
	"os"
	"strings"
)

//go:embed index.html
var indexHtmlTemplateContent string

//go:embed css/pico.min.css
var picoCssTemplateContent string

//go:embed js/main.js
var mainJsTemplateContent string

//go:embed js/htmx.min.js
var htmxJsTemplateContent string

//go:embed js/modal.js
var modalJsTemplateContent string

type TemplateFileType int

const (
	HtmlTemplateFileType TemplateFileType = 0
	CssTemplateFileType  TemplateFileType = 1
	JsTemplateFileType   TemplateFileType = 2
)

func GetIndexHtmlContent(frontendUri string, wsUri string) (string, error) {
	htmlTemplate, err := template.New("").Parse(indexHtmlTemplateContent)
	if err != nil {
		return "", fmt.Errorf("html.GetIndexHtmlContent: could not parse %v template: %w", "index.html", err)
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("html.GetIndexHtmlContent: could not get current working directory: %w", err)
	}

	availableVideosSelect, err := templ.GetAvailableVideosSelect(false)
	if err != nil {
		return "", fmt.Errorf("html.GetIndexHtmlContent: could not get available-videos-select component: %w", err)
	}
	availableVideosSelectBuffer := new(bytes.Buffer)
	availableVideosSelect.Render(context.Background(), availableVideosSelectBuffer)

	data := struct {
		HomeDirectory         string
		FrontendUri           string
		WsUri                 string
		AvailableVideosSelect template.HTML
	}{
		HomeDirectory:         currentDir,
		FrontendUri:           frontendUri,
		WsUri:                 wsUri,
		AvailableVideosSelect: template.HTML(availableVideosSelectBuffer.String()),
	}

	var templateBuffer bytes.Buffer
	err = htmlTemplate.Execute(&templateBuffer, data)
	if err != nil {
		return "", fmt.Errorf("html.GetIndexHtmlContent: could not execute %v template: %w", "index.html", err)
	}

	return templateBuffer.String(), nil
}

func GetMainJsContent(frontendUri string) (string, error) {
	templateData := struct {
		FrontendUri string
	}{
		FrontendUri: frontendUri,
	}

	content, err := getTemplateString(mainJsTemplateContent, JsTemplateFileType, templateData)
	if err != nil {
		return "", fmt.Errorf("html.GetMainJsContent: could not get js/main.js template string: %w", err)
	}

	content = strings.TrimPrefix(content, "<script>")
	content = strings.TrimSuffix(content, "</script>")

	return content, nil
}

func GetPicoCssContent() (string, error) {
	content, err := getTemplateString(picoCssTemplateContent, CssTemplateFileType, nil)
	if err != nil {
		return "", fmt.Errorf("html.GePicoCssContent: could not get css/pico.min.css template string: %w", err)
	}

	return content, nil
}

func GetHtmxJsContent() (string, error) {
	content, err := getTemplateString(htmxJsTemplateContent, JsTemplateFileType, nil)
	if err != nil {
		return "", fmt.Errorf("html.GetHtmxJsContent: could not get js/htmx.min.js template string: %w", err)
	}

	content = strings.TrimPrefix(content, "<script>")
	content = strings.TrimSuffix(content, "</script>")

	return content, nil
}

func GetModalJsContent() (string, error) {
	content, err := getTemplateString(modalJsTemplateContent, JsTemplateFileType, nil)
	if err != nil {
		return "", fmt.Errorf("html.GetModalJsContent: could not get js/modal.js template string: %w", err)
	}

	content = strings.TrimPrefix(content, "<script>")
	content = strings.TrimSuffix(content, "</script>")

	return content, nil
}

func getTemplateString(templateContent string, templateFileType TemplateFileType, data any) (string, error) {
	if templateFileType == JsTemplateFileType {
		templateContent = fmt.Sprintf("<script>%s</script>", templateContent)
	}

	template, err := template.New("").Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("html.getTemplateString: could not parse template: %w", err)
	}

	var templateBuffer bytes.Buffer
	err = template.Execute(&templateBuffer, data)
	if err != nil {
		return "", fmt.Errorf("html.getTemplateString: could not execute template: %w", err)
	}

	return templateBuffer.String(), nil
}
