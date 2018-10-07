package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
	"nekoserver/middleware/data"
	"nekoserver/middleware/func"
)

type administrator struct {
	Name string `json:"username"`
	Password string `json:"password"`
}

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

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestPath := c.Request.URL.Path
		if c.Request.Method == http.MethodPost {
			fmt.Println("Use Middleware")
			fmt.Printf("%s\n", requestPath)
			if !strings.Contains(requestPath, "5bb8745bb8436c5c964c04e03678364c5c964c2a5cc3b070") {
				c.Next()
			} else {
				tokenRequest, err := request.ParseFromRequest(c.Request, request.AuthorizationHeaderExtractor,
					func(token *jwt.Token) (interface{}, error) {
						return []byte("nekohandversion7.0"), nil
					})
				if err != nil {
					_func.RespondError(c, http.StatusUnauthorized,
						data.Error{
							Code: "401",
							Message: "未授权",
						})
					return
				}
				if tokenRequest.Valid {
					c.Next()
				} else {
					_func.RespondError(c, http.StatusUnauthorized,
						data.Error{
							Code: "401",
							Message: "未授权",
						})
					return
				}
			}
		} else {
			c.Next()
		}
	}
}