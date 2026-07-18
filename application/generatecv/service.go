package generatecv

import (
	"context"

	"cvgen/domain/resume"
)

// TemplateRenderer converte um Resume em HTML.
type TemplateRenderer interface {
	Render(r *resume.Resume) (string, error)
}

// PDFRenderer converte HTML em bytes de um PDF.
type PDFRenderer interface {
	RenderPDF(ctx context.Context, html string) ([]byte, error)
}

// UseCase orquestra a geração do CV: HTML -> PDF.
type UseCase struct {
	templateRenderer TemplateRenderer
	pdfRenderer      PDFRenderer
}

func NewUseCase(t TemplateRenderer, p PDFRenderer) *UseCase {
	return &UseCase{
		templateRenderer: t,
		pdfRenderer:      p,
	}
}

func (uc *UseCase) Execute(ctx context.Context, r *resume.Resume) ([]byte, error) {
	html, err := uc.templateRenderer.Render(r)
	if err != nil {
		return nil, err
	}

	pdfBytes, err := uc.pdfRenderer.RenderPDF(ctx, html)
	if err != nil {
		return nil, err
	}

	return pdfBytes, nil
}
