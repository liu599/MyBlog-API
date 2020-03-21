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

func TokenRemoteAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestPath := c.Request.URL.Path
		if !strings.Contains(requestPath, "auth") {
			c.Next()
			return
		}
		client := &http.Client{}
		form := url.Values{}
		requestUid, _ :=  c.Request.Header["User"]
		requestToken := c.Request.Header["Authorization"]
		form.Add("uid", requestUid[0])
		form.Add("token", requestToken[0])
		req, _ := http.NewRequest("POST", "http://localhost:19223/v2/auth/token.auth", strings.NewReader(form.Encode()))
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
			var p data.ResponseBody
			bodyBytes, _ := ioutil.ReadAll(resp.Body)
			json.Unmarshal(bodyBytes, &p)
			if !p.Valid {
				_func.RespondError(c, http.StatusBadRequest,
					data.Error{
						Code: strconv.Itoa(e.ERROR_AUTH_CHECK_TOKEN_FAIL),
						Message: fmt.Sprintf("%v", e.GetMsg(e.ERROR_AUTH_CHECK_TOKEN_FAIL)),
					})
				return
			}
		}
		c.Next()
	}
}