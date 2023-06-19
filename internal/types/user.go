package types

type User struct {
	ID    uint
	Name  string `json:"name"`
	Email string `json:"email"`
}
