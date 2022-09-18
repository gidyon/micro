package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

// Payload contains jwt payload
type Payload struct {
	ID           string
	ProjectID    string
	Names        string
	PhoneNumber  string
	EmailAddress string
	Group        string
	Roles        []string
}

// Claims contains JWT claims information
type Claims struct {
	*Payload
	jwt.StandardClaims
}

func (api *authAPI) genToken(ctx context.Context, payload *Payload, expires int64) (tokenStr string, err error) {
	defer func() {
		if err2 := recover(); err2 != nil {
			err = fmt.Errorf("%v", err2)
		}
	}()

	token := jwt.NewWithClaims(api.SigningMethod, Claims{
		Payload: payload,
		StandardClaims: jwt.StandardClaims{
			Audience:  api.Audience,
			ExpiresAt: expires,
			IssuedAt:  time.Now().Unix(),
			Issuer:    api.Issuer,
			NotBefore: 0,
			Subject:   "",
		},
	})

	token.Header["kid"] = payload.ProjectID

	return token.SignedString(api.SigningKey)
}

func (api *authAPI) genTokenV2(ctx context.Context, claims *Claims, expires int64, signingKey []byte) (tokenStr string, err error) {
	defer func() {
		if err2 := recover(); err2 != nil {
			err = fmt.Errorf("%v", err2)
		}
	}()

	token := jwt.NewWithClaims(api.SigningMethod, *claims)

	token.Header["kid"] = claims.ProjectID

	return token.SignedString(signingKey)
}
