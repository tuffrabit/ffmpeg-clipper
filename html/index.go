package html

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"os"
	"strings"
)

//go:embed index.html
var indexHtmlTemplateContent string

//go:embed index2.html
var index2HtmlTemplateContent string

//go:embed css/w3.css
var w3CssTemplateContent string

//go:embed css/w3-theme-blue-grey.css
var w3ThemeBlueGreyCssTemplateContent string

//go:embed css/pico.min.css
var picoCssTemplateContent string

//go:embed js/main.js
var mainJsTemplateContent string

//go:embed js/htmx.min.js
var htmxJsTemplateContent string

//go:embed js/modal.js
var modalJsTemplateContent string

//go:embed encoder/libx264.html
var libx264EncoderFieldsTemplateContent string

//go:embed encoder/libx265.html
var libx265EncoderFieldsTemplateContent string

//go:embed encoder/libaomav1.html
var libaomAv1EncoderFieldsTemplateContent string

//go:embed encoder/nvench264.html
var nvencH264EncoderFieldsTemplateContent string

//go:embed encoder/nvenchevc.html
var nvencHevcEncoderFieldsTemplateContent string

//go:embed encoder/intelh264.html
var intelH264EncoderFieldsTemplateContent string

//go:embed encoder/intelhevc.html
var intelHevcEncoderFieldsTemplateContent string

//go:embed encoder/intelav1.html
var intelAv1EncoderFieldsTemplateContent string

type TemplateFileType int

const (
	HtmlTemplateFileType TemplateFileType = 0
	CssTemplateFileType  TemplateFileType = 1
	JsTemplateFileType   TemplateFileType = 2
)

func GetIndexHtmlContent(frontendUri string, wsUri string) (string, error) {
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
		WsUri         string
	}{
		HomeDirectory: currentDir,
		FrontendUri:   frontendUri,
		WsUri:         wsUri,
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

	intelH264EncoderFieldsContent, err := getTemplateString(intelH264EncoderFieldsTemplateContent, HtmlTemplateFileType, nil)
	if err != nil {
		return "", fmt.Errorf("html.GetIndex2HtmlContent: could not get encoder/intelh264.html template string: %w", err)
	}

	intelHevcEncoderFieldsContent, err := getTemplateString(intelHevcEncoderFieldsTemplateContent, HtmlTemplateFileType, nil)
	if err != nil {
		return "", fmt.Errorf("html.GetIndex2HtmlContent: could not get encoder/intelhevc.html template string: %w", err)
	}

	intelAv1EncoderFieldsContent, err := getTemplateString(intelAv1EncoderFieldsTemplateContent, HtmlTemplateFileType, nil)
	if err != nil {
		return "", fmt.Errorf("html.GetIndex2HtmlContent: could not get encoder/intelav1.html template string: %w", err)
	}

	indexHtmlTemplateData := struct {
		HomeDirectory              string
		FrontendUri                string
		W3Css                      template.HTML
		W3ThemeBlueGreyCss         template.HTML
		Libx264EncoderFieldsHtml   template.HTML
		Libx265EncoderFieldsHtml   template.HTML
		LibaomAv1EncoderFieldsHtml template.HTML
		NvencH264EncoderFieldHtml  template.HTML
		NvencHevcEncoderFieldHtml  template.HTML
		IntelH264EncoderFieldHtml  template.HTML
		IntelHevcEncoderFieldHtml  template.HTML
		IntelAv1EncoderFieldHtml   template.HTML
	}{
		HomeDirectory:              currentDir,
		FrontendUri:                frontendUri,
		W3Css:                      template.HTML(w3CssContent),
		W3ThemeBlueGreyCss:         template.HTML(w3ThemeBlueGreyCssContent),
		Libx264EncoderFieldsHtml:   template.HTML(libx264EncoderFieldsContent),
		Libx265EncoderFieldsHtml:   template.HTML(libx265EncoderFieldsContent),
		LibaomAv1EncoderFieldsHtml: template.HTML(libaomAv1EncoderFieldsContent),
		NvencH264EncoderFieldHtml:  template.HTML(nvencH264EncoderFieldsContent),
		NvencHevcEncoderFieldHtml:  template.HTML(nvencHevcEncoderFieldsContent),
		IntelH264EncoderFieldHtml:  template.HTML(intelH264EncoderFieldsContent),
		IntelHevcEncoderFieldHtml:  template.HTML(intelHevcEncoderFieldsContent),
		IntelAv1EncoderFieldHtml:   template.HTML(intelAv1EncoderFieldsContent),
	}

	indexHtml, err := getTemplateString(index2HtmlTemplateContent, HtmlTemplateFileType, indexHtmlTemplateData)
	if err != nil {
		return "", fmt.Errorf("html.GetIndex2HtmlContent: could not get index2.html template string: %w", err)
	}

	return indexHtml, nil
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
