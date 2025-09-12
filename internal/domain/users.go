package domain

type User struct {
	ID       string `json:"id"`       // Student/Teacher ID
	Username string `json:"username"` // FIO Student
	Role     string `json:"role"`     // Teacher, Admin, People (Students)
}
