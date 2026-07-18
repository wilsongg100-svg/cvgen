package importcv

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"

	"cvgen/domain/resume"
)

const (
	markerStart = "CVGEN_DATA_START:"
	markerEnd   = ":CVGEN_DATA_END"
)

// TextExtractor extrai o texto puro de um ficheiro PDF.
type TextExtractor interface {
	Extract(pdfBytes []byte) (string, error)
}

// UseCase reconstrói um Resume a partir de um PDF gerado por esta própria API.
type UseCase struct {
	textExtractor TextExtractor
}

func NewUseCase(t TextExtractor) *UseCase {
	return &UseCase{textExtractor: t}
}

func (uc *UseCase) Execute(pdfBytes []byte) (*resume.Resume, error) {
	text, err := uc.textExtractor.Extract(pdfBytes)
	if err != nil {
		return nil, err
	}

	startIdx := strings.Index(text, markerStart)
	endIdx := strings.Index(text, markerEnd)
	if startIdx == -1 || endIdx == -1 || endIdx < startIdx {
		return nil, errors.New("este PDF não contém os dados gerados por esta API (ou foi editado/reimpresso)")
	}

	encoded := text[startIdx+len(markerStart) : endIdx]
	// a extração de PDF pode introduzir espaços/quebras de linha no meio do base64
	encoded = strings.Join(strings.Fields(encoded), "")

	raw, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, errors.New("dados embutidos corrompidos ou inválidos")
	}

	var r resume.Resume
	if err := json.Unmarshal(raw, &r); err != nil {
		return nil, errors.New("falha ao interpretar os dados do CV")
	}

	return &r, nil
}
