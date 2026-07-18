package resume

import "errors"

type Candidate struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone,omitempty"`
	Location string `json:"location,omitempty"`
	LinkedIn string `json:"linkedin,omitempty"`
	GitHub   string `json:"github,omitempty"`
	Website  string `json:"website,omitempty"`
	Summary  string `json:"summary,omitempty"`
}

type Experience struct {
	Role         string   `json:"role"`
	Company      string   `json:"company"`
	Location     string   `json:"location,omitempty"`
	StartDate    string   `json:"startDate"`
	EndDate      string   `json:"endDate,omitempty"`
	Achievements []string `json:"achievements,omitempty"`
}

type Education struct {
	Degree      string `json:"degree"`
	Institution string `json:"institution"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
	Grade       string `json:"grade,omitempty"`
}

type SkillGroup struct {
	Label string   `json:"label"`
	Items []string `json:"items"`
}

type Resume struct {
	Candidate   Candidate    `json:"candidate"`
	Experience  []Experience `json:"experience,omitempty"`
	Education   []Education  `json:"education,omitempty"`
	SkillGroups []SkillGroup `json:"skillGroups,omitempty"`
	Languages   []string     `json:"languages,omitempty"`
}

func NewResume(
	candidate Candidate,
	experience []Experience,
	education []Education,
	skillGroups []SkillGroup,
	languages []string,
) (*Resume, error) {
	if candidate.Name == "" {
		return nil, errors.New("candidate name is required")
	}
	if candidate.Email == "" {
		return nil, errors.New("candidate email is required")
	}
	if len(experience) == 0 && len(education) == 0 {
		return nil, errors.New("resume must have at least one experience or education entry")
	}

	return &Resume{
		Candidate:   candidate,
		Experience:  experience,
		Education:   education,
		SkillGroups: skillGroups,
		Languages:   languages,
	}, nil
}
