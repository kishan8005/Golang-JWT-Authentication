package helper

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type SignedDetails struct {
	Email      string
	First_name string
	Last_name  string
	Uid        string
	User_type  string
	jwt.StandardClaims
}

var secretKey = os.Getenv("SECRET_KEY")

func GenerateAllTokens(email, fname, lname, phno, userId, userType string) (token string, refresh_token string, err error) {
	claims := &SignedDetails{
		Email:      email,
		First_name: fname,
		Uid:        userId,
		User_type:  userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	token, err = jwt.NewWithClaims(jwt.SigningMethodES256, claims).SignedString(secretKey)
	refresh_token, err = jwt.NewWithClaims(jwt.SigningMethodES256, refreshClaims).SignedString(secretKey)

	if err != nil {
		log.Panic(err)
		return
	}

	return token, refresh_token, err
}
