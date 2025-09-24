package domain

type User struct {
	ID       string `json:"id"`       // Student/Teacher ID
	Username string `json:"username"` // FIO Student
	Role     string `json:"role"`     // Teacher, Admin, People (Students)
}

type UserExtended struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Role          string `json:"role"`
	AcademicGroup string `json:"academic_group"`
	Profile       string `json:"profile,omitempty"`
}
