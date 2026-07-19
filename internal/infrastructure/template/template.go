package template

import (
	"bytes"
	"cvgen/internal/resume"
	"encoding/base64"
	"encoding/json"
	htmltemplate "html/template"
)

type viewData struct {
	*resume.Resume
	EmbeddedData string
}

type Engine struct {
	tmpl *htmltemplate.Template
}

func NewEngine(templatePath string) (*Engine, error) {
	tmpl, err := htmltemplate.ParseFiles(templatePath)
	if err != nil {
		return nil, err
	}
	return &Engine{tmpl: tmpl}, nil
}

func (e *Engine) Render(r *resume.Resume) (string, error) {
	raw, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	encoded := base64.StdEncoding.EncodeToString(raw)

	data := viewData{
		Resume:       r,
		EmbeddedData: encoded,
	}

	var buf bytes.Buffer
	if err := e.tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
