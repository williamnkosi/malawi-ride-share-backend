package models

type CreateUser struct {
	FirstName   string
	LastName    string
	PhoneNumber string
	Email       string
	Password1   string
	Password2   string
	Role        string
}

const (
	passenger = "passenger"
	driver    = "driver"
	admin     = "admin"
	support   = "support"
)
