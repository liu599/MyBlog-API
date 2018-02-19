package _func

import (
	"fmt"

	"nekoserver/middleware/data"

	"github.com/gin-gonic/gin"
)


func Respond(context *gin.Context, code int, data ...map[string]interface{}) {

	res := gin.H{}

	for _, v := range data {
		fmt.Println(v)
		for k, m := range v {
			switch t:=m.(type) {
			default:
				res[k] = t
			}
		}
	}

	res["success"] = true
	res["code"] = 0

	Response(context, code, res)

	context.Abort()
}

func RespondError(context *gin.Context, code int, err data.Error) {

	data := gin.H{}

	data["success"] = false
	data["error"] = err
	data["code"] = 1

	Response(context, code, data)

	context.Abort()
}

// Response
func Response(context *gin.Context, code int, data gin.H) {
	context.JSON(code, data)
}