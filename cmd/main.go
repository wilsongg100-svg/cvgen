package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"cvgen/application/generatecv"
	"cvgen/application/importcv"
	"cvgen/domain/resume"
	"cvgen/infrastructure/pdftext"
	"cvgen/infrastructure/renderer"
	tmpl "cvgen/infrastructure/template"
)

func main() {
	engine, err := tmpl.NewEngine("templates/cv.html")
	if err != nil {
		log.Fatalf("failed to load template: %v", err)
	}

	pdfRenderer := renderer.NewChromedpRenderer()
	useCase := generatecv.NewUseCase(engine, pdfRenderer)
	textExtractor := pdftext.NewExtractor()
	importUseCase := importcv.NewUseCase(textExtractor)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/resumes", handleGenerateCV(useCase))
	mux.HandleFunc("POST /api/v1/resumes/import", handleImportCV(importUseCase))

	log.Println("server running on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

func handleGenerateCV(uc *generatecv.UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var dto resume.RequestDTO
		if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
			http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		candidate := resume.Candidate{
			Name:     dto.Name,
			Email:    dto.Email,
			Phone:    dto.Phone,
			Location: dto.Location,
			LinkedIn: dto.LinkedIn,
			GitHub:   dto.GitHub,
			Website:  dto.Website,
			Summary:  dto.Summary,
		}

		var experience []resume.Experience
		for _, e := range dto.Experience {
			experience = append(experience, resume.Experience{
				Role:         e.Role,
				Company:      e.Company,
				StartDate:    e.StartDate,
				EndDate:      e.EndDate,
				Achievements: e.Achievements,
				Location:     e.Location,
			})
		}
		var education []resume.Education
		for _, e := range dto.Education {
			education = append(education, resume.Education{
				Degree:      e.Degree,
				Institution: e.Institution,
				StartDate:   e.StartDate,
				EndDate:     e.EndDate,
			})
		}

		var skillGroups []resume.SkillGroup
		for _, s := range dto.SkillGroups {
			skillGroups = append(skillGroups, resume.SkillGroup{
				Label: s.Label,
				Items: s.Items,
			})
		}

		cv, err := resume.NewResume(candidate, experience, education, skillGroups, dto.Languages)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
		defer cancel()

		pdfBytes, err := uc.Execute(ctx, cv)
		if err != nil {
			http.Error(w, "failed to generate pdf: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", "attachment; filename=cv.pdf")
		w.Write(pdfBytes)
	}
}

func handleImportCV(uc *importcv.UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file, _, err := r.FormFile("cv")
		if err != nil {
			http.Error(w, "envia o ficheiro no campo 'cv' (multipart/form-data)", http.StatusBadRequest)
			return
		}
		defer file.Close()

		pdfBytes, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "falha ao ler o ficheiro", http.StatusBadRequest)
			return
		}

		cv, err := uc.Execute(pdfBytes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cv)
	}
}
