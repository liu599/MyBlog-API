package controller

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func generateToken(uid int, environmentVariable string) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
	claims["iat"] = time.Now().Unix()
	token.Claims = claims
	fmt.Printf("\n%v\n", uid)
	// str1 := strconv.Itoa(uid)
	tokenString, err := token.SignedString([]byte(environmentVariable))

	fmt.Printf("\n%v\n", tokenString)

	if err != nil {
		fmt.Println("error in convert tokenString")
		return "Cannot convert secret string", err
	}

	return tokenString, nil
}

