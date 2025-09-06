package domain

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Mail     string `json:"mail"`
	Role     string `json:"role"`
}
