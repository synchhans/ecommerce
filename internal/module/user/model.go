package user

type User struct {
	ID     string `json:"id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type AuthResult struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}
