package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"gopkg.in/mgo.v2/bson"
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

func PostFetchOne(context *gin.Context) {
	id := context.Param("pid")
	err, post := models.PostFetchOne(id)
	if err != nil {
		_func.Respond(context, http.StatusBadRequest, gin.H{"error": err})
		return
	}
	mk := make(map[string]interface{})
	mk["data"] = post
	_func.Respond(context, http.StatusOK, mk)
}

func PostsChornology(context *gin.Context) {
	err, chr := models.PostsFetchChronology()
	if err != nil {
		_func.Respond(context, http.StatusBadRequest, gin.H{"error": err})
		return
	}
	mk := make(map[string]interface{})
	mk["data"] = chr
	_func.Respond(context, http.StatusOK, mk)
}

func PostCreation(context *gin.Context) {

	var p data.Post

	decoder := json.NewDecoder(context.Request.Body)
	if err := decoder.Decode(&p); err != nil {
		_func.RespondError(context, http.StatusBadRequest, data.Error{
			Message: "Invalid Request Payload",
		})
		return
	}

	defer context.Request.Body.Close()

	err, postId := models.CreatePost(p)

	if err != nil {
		_func.RespondError(context, http.StatusInternalServerError, data.Error{
			Message: "Database Error, Fail to create the post",
		})
		return
	}
	mk := make(map[string]interface{})
	mk["data"] = "a post"+ postId + " has been successful created " + bson.NewObjectId().Hex()
	_func.Respond(context, http.StatusOK, mk)
}