package domain

type User struct {
	ID       string `json:"id" bson:"_id"`
	Username string `json:"username" bson:"username"`
	Role     string `json:"role" bson:"role"`
}
