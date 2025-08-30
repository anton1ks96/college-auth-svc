package domain

type Student struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Group    string `json:"group"`
	Role     string `json:"role"`
}
