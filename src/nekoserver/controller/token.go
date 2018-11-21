package controller

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/json-iterator/go"
	"golang.org/x/crypto/scrypt"
	"nekoserver/e"
	"nekoserver/middleware/auth"
	"nekoserver/middleware/data"
	"nekoserver/middleware/func"
	"nekoserver/models"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func TokenFetch(context *gin.Context) {
	user := context.PostForm("username")
	password := context.PostForm("password")
	client := &http.Client{}
	form := url.Values{}
	form.Add("username", user)
	form.Add("password", password)
	req, _ := http.NewRequest("POST", "http://localhost:19223/v2/auth/token.get", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, _ := client.Do(req)
	if resp.StatusCode == http.StatusBadRequest {
		_func.RespondError(context, http.StatusBadRequest,
			data.Error{
				Code: strconv.Itoa(e.ERROR_AUTH_CHECK_TOKEN_FAIL),
				Message: fmt.Sprintf("%v", e.GetMsg(e.ERROR_AUTH_CHECK_TOKEN_FAIL)),
			})
		return
	}
	if resp.StatusCode == http.StatusOK {
		var p data.ResponseBody
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(bodyBytes, &p)
		mk := make(map[string]interface{})
		mk["api_token"] = p.Token
		mk["user_id"] = p.Uid
		_func.Respond(context, http.StatusOK, mk)
	}
}

func TokenGen(context *gin.Context) {
	var usr data.User
	user := context.PostForm("username")
	password := context.PostForm("password")
	usr.Name = user
	dk, _ := scrypt.Key([]byte(password), []byte(os.Getenv("PASS_GEN")), 16384, 8, 1, 32)
	usr.Password = base64.StdEncoding.EncodeToString(dk)
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

