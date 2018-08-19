package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"nekoserver/middleware/data"
	"nekoserver/middleware/func"
	"nekoserver/models"

	"github.com/gin-gonic/gin"
)

func CategoriesFetch(context *gin.Context) {

	err, cats := models.FetchCategoryList()

	mk := make(map[string]interface{})

	mk["data"] = cats

	if err != nil {
		fmt.Println(err)
		_func.RespondError(context, http.StatusInternalServerError, data.Error{
			Code: strconv.Itoa(http.StatusInternalServerError),
			Message: "服务器内部错误",
		})
		return
	}

	_func.Respond(context, http.StatusOK, mk)

}