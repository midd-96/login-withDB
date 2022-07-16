package helpers

import (
	"log"
	"os"
	"project_login/models"
	"time"

	"github.com/golang-jwt/jwt"
)

var SECRET_KEY string = os.Getenv("SECRET_KEY")

func GenerateTokens(username string, usertype string) (signedToken string, signedRefreshToken string, err error) {

	claims := &models.SignedDetails{
		Username:  username,
		User_type: usertype,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}
	refreshClaims := &models.SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
		return
	}
	return token, refreshToken, err
}

func ValidateToken(signedToken string) bool {
	check := true
	token, err := jwt.ParseWithClaims(
		signedToken,
		&models.SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	claims, ok := token.Claims.(*models.SignedDetails)

	if !ok {
		log.Println("The token is invalid")
		check = false
		return check
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		log.Println("Token Expired")
		check = false
		return check
	}
	return check

}
