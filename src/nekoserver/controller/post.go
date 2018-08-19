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

func PostsFetch(context *gin.Context) {
	pageNumber, _ := strconv.Atoi(context.PostForm("pageNumber"))
	pageSize, _ := strconv.Atoi(context.PostForm("pageSize"))
	fmt.Println(pageSize, pageNumber)
	err, posts := models.PostsFetchAllWithPageNumber((pageNumber - 1) * pageSize, pageSize)
	if err != nil {
		_func.Respond(context, http.StatusBadRequest, gin.H{"error": err})
		return
	}
	mk := make(map[string]interface{})

	mk["data"] = posts

	_, totalNumber := models.PostsFetchTotalNumber()
	mk["pager"]  = data.Pager{
		PageNum: pageNumber,
		PageSize: pageSize,
		TotalNumber: totalNumber,
	}
	_func.Respond(context, http.StatusOK, mk)
}

func PostsFetchByCategory(context *gin.Context) {
	pageNumber, _ := strconv.Atoi(context.PostForm("pageNumber"))
	pageSize, _ := strconv.Atoi(context.PostForm("pageSize"))
	cid := context.Param("cid")
	err, posts := models.PostsFetchCategoryWithPageNumber((pageNumber - 1) * pageSize, pageSize, cid)

	if err != nil {
		_func.Respond(context, http.StatusBadRequest, gin.H{"error": err})
		return
	}
	mk := make(map[string]interface{})

	mk["data"] = posts

	_, totalCategoryNumber := models.PostsFetchTotalNumberByCategory(cid)
	mk["pager"]  = data.Pager{
		PageNum: pageNumber,
		PageSize: pageSize,
		TotalNumber: totalCategoryNumber,
	}

	_func.Respond(context, http.StatusOK, mk)
}