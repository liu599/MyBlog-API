package _func

import (

	"nekoserver/middleware/data"

	"github.com/gin-gonic/gin"
)


func Respond(context *gin.Context, code int, data ...map[string]interface{}) {

	res := gin.H{}

	for _, v := range data {
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

	emptyData := gin.H{}

	emptyData["success"] = false
	emptyData["error"] = err
	emptyData["code"] = 1

	Response(context, code, emptyData)

	context.Abort()
}

// Response
func Response(context *gin.Context, code int, data gin.H) {
	context.JSON(code, data)
}