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

//go:embed index2.html
var index2HtmlTemplateContent string

//go:embed css/w3.css
var w3CssTemplateContent string

//go:embed css/w3-theme-blue-grey.css
var w3ThemeBlueGreyCssTemplateContent string

//go:embed js/main.js
var mainJsTemplateContent string

//go:embed encoder/libx264.html
var libx264EncoderFieldsTemplateContent string

//go:embed encoder/libx265.html
var libx265EncoderFieldsTemplateContent string

//go:embed encoder/nvench264.html
var nvencH264EncoderFieldsTemplateContent string

//go:embed encoder/nvenchevc.html
var nvencHevcEncoderFieldsTemplateContent string

//go:embed encoder/libaomav1.html
var libaomAv1EncoderFieldsTemplateContent string

type TemplateFileType int

const (
	HtmlTemplateFileType TemplateFileType = 0
	CssTemplateFileType  TemplateFileType = 1
	JsTemplateFileType   TemplateFileType = 2
)

func GetIndexHtmlContent(frontendUri string) (string, error) {
	template, err := template.New("").Parse(indexHtmlTemplateContent)
	if err != nil {
		return "", fmt.Errorf("html.GetIndexHtmlContent: could not parse %v template: %w", "index.html", err)
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("html.GetIndexHtmlContent: could not get current working directory: %w", err)
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
		return "", fmt.Errorf("html.GetIndexHtmlContent: could not execute %v template: %w", "index.html", err)
	}

	return templateBuffer.String(), nil
}

func GetIndex2HtmlContent(frontendUri string) (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("html.GetIndex2HtmlContent: could not get current working directory: %w", err)
	}

	w3CssContent, err := getTemplateString(w3CssTemplateContent, CssTemplateFileType, nil)
	if err != nil {
		return "", fmt.Errorf("html.GetIndex2HtmlContent: could not get css/w3.css template string: %w", err)
	}

	w3ThemeBlueGreyCssContent, err := getTemplateString(w3ThemeBlueGreyCssTemplateContent, CssTemplateFileType, nil)
	if err != nil {
		return "", fmt.Errorf("html.GetIndex2HtmlContent: could not get css/w3-theme-blue-grey.css template string: %w", err)
	}

	mainJsContent, err := getTemplateString(mainJsTemplateContent, JsTemplateFileType, nil)
	if err != nil {
		return "", fmt.Errorf("html.GetIndex2HtmlContent: could not get js/main.js template string: %w", err)
	}

	libx264EncoderFieldsContent, err := getTemplateString(libx264EncoderFieldsTemplateContent, HtmlTemplateFileType, nil)
	if err != nil {
		return "", fmt.Errorf("html.GetIndex2HtmlContent: could not get encoder/libx264.html template string: %w", err)
	}

	libx265EncoderFieldsContent, err := getTemplateString(libx265EncoderFieldsTemplateContent, HtmlTemplateFileType, nil)
	if err != nil {
		return "", fmt.Errorf("html.GetIndex2HtmlContent: could not get encoder/libx265.html template string: %w", err)
	}

	libaomAv1EncoderFieldsContent, err := getTemplateString(libaomAv1EncoderFieldsTemplateContent, HtmlTemplateFileType, nil)
	if err != nil {
		return "", fmt.Errorf("html.GetIndex2HtmlContent: could not get encoder/libaomav1.html template string: %w", err)
	}

	nvencH264EncoderFieldsContent, err := getTemplateString(nvencH264EncoderFieldsTemplateContent, HtmlTemplateFileType, nil)
	if err != nil {
		return "", fmt.Errorf("html.GetIndex2HtmlContent: could not get encoder/nvench264.html template string: %w", err)
	}

	nvencHevcEncoderFieldsContent, err := getTemplateString(nvencHevcEncoderFieldsTemplateContent, HtmlTemplateFileType, nil)
	if err != nil {
		return "", fmt.Errorf("html.GetIndex2HtmlContent: could not get encoder/nvenchevc.html template string: %w", err)
	}

	indexHtmlTemplateData := struct {
		HomeDirectory              string
		FrontendUri                string
		W3Css                      template.HTML
		W3ThemeBlueGreyCss         template.HTML
		MainJs                     template.HTML
		Libx264EncoderFieldsHtml   template.HTML
		Libx265EncoderFieldsHtml   template.HTML
		LibaomAv1EncoderFieldsHtml template.HTML
		NvencH264EncoderFieldHtml  template.HTML
		NvencHevcEncoderFieldHtml  template.HTML
	}{
		HomeDirectory:              currentDir,
		FrontendUri:                frontendUri,
		W3Css:                      template.HTML(w3CssContent),
		W3ThemeBlueGreyCss:         template.HTML(w3ThemeBlueGreyCssContent),
		MainJs:                     template.HTML(mainJsContent),
		Libx264EncoderFieldsHtml:   template.HTML(libx264EncoderFieldsContent),
		Libx265EncoderFieldsHtml:   template.HTML(libx265EncoderFieldsContent),
		LibaomAv1EncoderFieldsHtml: template.HTML(libaomAv1EncoderFieldsContent),
		NvencH264EncoderFieldHtml:  template.HTML(nvencH264EncoderFieldsContent),
		NvencHevcEncoderFieldHtml:  template.HTML(nvencHevcEncoderFieldsContent),
	}

	indexHtml, err := getTemplateString(index2HtmlTemplateContent, HtmlTemplateFileType, indexHtmlTemplateData)
	if err != nil {
		return "", fmt.Errorf("html.GetIndex2HtmlContent: could not get index2.html template string: %w", err)
	}

	return indexHtml, nil
}

func getTemplateString(templateContent string, templateFileType TemplateFileType, data any) (string, error) {
	if templateFileType == CssTemplateFileType {
		templateContent = fmt.Sprintf("<style>%s</style>", templateContent)
	} else if templateFileType == JsTemplateFileType {
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
