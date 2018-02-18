package _func

import "github.com/gin-gonic/gin"

func Respond(context *gin.Context, code int, data ...interface{}) {
	var res = make(map[string]interface{})

	res["data"] = data
	res["success"] = true
	res["code"] = 0

	Response(context, code, res)
}

// Response
// 响应
func Response(context *gin.Context, code int, data gin.H) {
	context.JSON(code, data)
}