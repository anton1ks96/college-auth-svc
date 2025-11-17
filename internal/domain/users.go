package domain

type User struct {
	ID       string `json:"id"`       // Student/Teacher ID
	Username string `json:"username"` // FIO Student
	Role     string `json:"role"`     // Teacher, Admin, People (Students)
}

type UserGroups struct {
	AcademicGroup string `json:"academic_group"`
	Profile       string `json:"profile,omitempty"`
	Subgroup      string `json:"subgroup,omitempty"`
	EnglishGroup  string `json:"english_group,omitempty"`
}

type UserExtended struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Role          string `json:"role"`
	AcademicGroup string `json:"academic_group"`
	Profile       string `json:"profile,omitempty"`
	Subgroup      string `json:"subgroup,omitempty"`
	EnglishGroup  string `json:"english_group,omitempty"`
}
