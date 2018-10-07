package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"nekoserver/middleware/data"
	"nekoserver/middleware/func"
	"nekoserver/models"
)

func CommentsFetch(context *gin.Context) {
	id := context.Param("pid")
	fmt.Println(id)
	err, comments := models.CommentsFetch(id)
	if err != nil {
		_func.Respond(context, http.StatusBadRequest, gin.H{"error": err})
		return
	}
	mk := make(map[string]interface{})

	mk["data"] = comments

	_func.Respond(context, http.StatusOK, mk)
}

func CommentCreation(context *gin.Context) {

	var co data.Comment

	decoder := json.NewDecoder(context.Request.Body)
	if err := decoder.Decode(&co); err != nil {
		fmt.Sprintf("%v", err)
		_func.RespondError(context, http.StatusBadRequest, data.Error{})
		return
	}

	defer context.Request.Body.Close()

	if err2 := models.CommentCreate(co); err2 != nil {
		_func.RespondError(context, http.StatusBadRequest, data.Error{})
		return
	}

	mk := make(map[string]interface{})
	mk["status"] = "successful created!"
	_func.Respond(context, http.StatusOK, mk)
}