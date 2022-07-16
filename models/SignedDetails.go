package models

import "github.com/golang-jwt/jwt"

type SignedDetails struct {
	Username  string
	User_type string
	jwt.StandardClaims
}
