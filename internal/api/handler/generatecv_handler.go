package handler

import (
	"context"
	"cvgen/internal/application/generatecv"
	"cvgen/internal/resume"
	"encoding/json"
	"net/http"
	"time"
)

type GeneratecvHandler struct {
	service *generatecv.UseCase
}

func NewGeneratecvHandler(service *generatecv.UseCase) *GeneratecvHandler {
	return &GeneratecvHandler{
		service: service,
	}
}

func (g *GeneratecvHandler) Generatecv(w http.ResponseWriter, r *http.Request) {
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

	pdfBytes, err := g.service.Execute(ctx, cv)
	if err != nil {
		http.Error(w, "failed to generate pdf: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=cv.pdf")
	_, err = w.Write(pdfBytes)
	if err != nil {
		http.Error(w, "failed to generate pdf: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
