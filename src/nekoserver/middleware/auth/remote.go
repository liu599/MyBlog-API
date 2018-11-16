package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"nekoserver/e"
	"nekoserver/middleware/data"
	"nekoserver/middleware/func"
)

type ResponseBody struct{
	Code int
	Success bool
	Token string
}


func TokenRemoteAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		client := &http.Client{}
		form := url.Values{}
		form.Add("username", "tokei")
		form.Add("password", "!7d4a3eEDDIE")
		req, _ := http.NewRequest("POST", "http://localhost:19223/v2/auth/token.get", strings.NewReader(form.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		resp, _ := client.Do(req)
		fmt.Println(resp.StatusCode)
		if resp.StatusCode == http.StatusBadRequest {
			_func.RespondError(c, http.StatusBadRequest,
				data.Error{
					Code: strconv.Itoa(e.ERROR_AUTH_CHECK_TOKEN_FAIL),
					Message: fmt.Sprintf("%v", e.GetMsg(e.ERROR_AUTH_CHECK_TOKEN_FAIL)),
				})
			return
		}
		if resp.StatusCode == http.StatusOK {
			var p ResponseBody
			bodyBytes, _ := ioutil.ReadAll(resp.Body)
			json.Unmarshal(bodyBytes, &p)
			fmt.Println(p.Token)
		}
		c.Next()
	}
}