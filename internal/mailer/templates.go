package mailer

import (
	"bytes"
	_ "embed"
	"fmt"
	"text/template"
)

//go:embed templates/verification_code_template.html
var verificationCodeTemplate string

func RenderVerificationCodeTemplate(code string) string {
	tmpl, _ := RenderTemplate("verification_code", verificationCodeTemplate, map[string]string{
		"Code": code,
	})

	return tmpl
}

func RenderTemplate(name, text string, data any) (string, error) {

	var buffer bytes.Buffer

	tmpl, err := template.New(name).Parse(text)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	if err := tmpl.Execute(&buffer, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buffer.String(), nil
}
