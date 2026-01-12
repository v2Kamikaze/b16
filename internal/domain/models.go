package domain

type User struct {
	ID       string
	Email    string
	Password string
	Roles    []Role
}

type Role struct {
	ID   string
	Name string
}
