package controller

import (
	"fmt"
	"net/http"

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
		_func.Respond(context, http.StatusBadRequest, nil)
	}

	_func.Respond(context, http.StatusOK, mk)

}