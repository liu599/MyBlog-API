package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"nekoserver/middleware/func"
	"nekoserver/models"

	"github.com/gin-gonic/gin"
)

func PostsFetch(context *gin.Context) {
	pageNumber, _ := strconv.Atoi(context.PostForm("pageNumber"))
	pageSize, _ := strconv.Atoi(context.PostForm("pageSize"))
	fmt.Println(pageSize, pageNumber)
	err, posts := models.PostsFetchAllWithPageNumber(0, 10)
	if err != nil {
		_func.Respond(context, http.StatusBadRequest, gin.H{"error": err})
		return
	}
	mk := make(map[string]interface{})

	mk["data"] = posts
	_func.Respond(context, http.StatusOK, mk)
}