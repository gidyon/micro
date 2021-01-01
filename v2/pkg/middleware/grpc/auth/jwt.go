package auth

import (
	"github.com/dgrijalva/jwt-go"
)

// Payload contains jwt payload
type Payload struct {
	ID           string
	ProjectID    string
	Names        string
	PhoneNumber  string
	EmailAddress string
	Group        string
}

// Claims contains JWT claims information
type Claims struct {
	*Payload
	jwt.StandardClaims
}
