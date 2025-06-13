package models

type User struct {
	ID       string
	Username string
	Email    string
}

type AuthResponse struct {
	Token string
	User  User
}

type Message struct{}

type Channel struct{}
