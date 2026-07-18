package resume

type experienceDTO struct {
	Role         string   `json:"role"`
	Company      string   `json:"company"`
	StartDate    string   `json:"startDate"`
	EndDate      string   `json:"endDate"`
	Location     string   `json:"location"`
	Achievements []string `json:"achievements"`
}

type educationDTO struct {
	Degree      string `json:"degree"`
	Institution string `json:"institution"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
}

type RequestDTO struct {
	Name        string          `json:"name"`
	Email       string          `json:"email"`
	Phone       string          `json:"phone"`
	Location    string          `json:"location"`
	LinkedIn    string          `json:"linkedin"`
	GitHub      string          `json:"github"`
	Website     string          `json:"website"`
	Summary     string          `json:"summary"`
	Experience  []experienceDTO `json:"experience"`
	Education   []educationDTO  `json:"education"`
	Skills      []string        `json:"skills"`
	SkillGroups []SkillGroup    `json:"skillgroups"`
	Languages   []string        `json:"languages"`
}
