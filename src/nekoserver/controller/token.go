package controller

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/scrypt"
	"nekoserver/middleware/auth"
	"nekoserver/middleware/data"
	"nekoserver/middleware/func"
	"nekoserver/models"
)

func TokenGen(context *gin.Context) {
	var usr data.User
	user := context.PostForm("username")
	password := context.PostForm("password")
	usr.Name = user
	dk, _ := scrypt.Key([]byte(password), []byte(os.Getenv("PASS_GEN")), 16384, 8, 1, 32)
	usr.Password = base64.StdEncoding.EncodeToString(dk)
	fmt.Println(password)
	fmt.Println(os.Getenv("PASS_GEN"))
	fmt.Println(dk)
	mk := make(map[string]interface{})
	err, sign := models.TokenCheckUser(usr)
	if err != nil {
		_func.RespondError(context, http.StatusInternalServerError, data.Error{
			Message: fmt.Sprintf("%v", err),
		})
		return
	}
	if sign == true {
		salt := os.Getenv("NEKO_TOKEN")
		token, _ := auth.GenerateToken(usr.Name, salt)
		mk["api_token"] = token
	} else {
		_func.RespondError(context, http.StatusUnauthorized, data.Error{
			Message: "Fail to validate this token. Please check your username and password",
		})
		return
	}
	_func.Respond(context, http.StatusOK, mk)
}

