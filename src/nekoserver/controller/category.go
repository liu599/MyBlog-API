package controller
<<<<<<< HEAD
=======

import (
	"fmt"
	"net/http"
	"strconv"

	"gopkg.in/mgo.v2/bson"
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

func CategoriesEdit(context *gin.Context) {
	var p data.Category

	decoder := json.NewDecoder(context.Request.Body)
	if err := decoder.Decode(&p); err != nil {
		_func.RespondError(context, http.StatusBadRequest, data.Error{
			Message: "Invalid Request Payload",
		})
		return
	}

	defer context.Request.Body.Close()

	err, pflag := models.FindCategory(p)

	if pflag == true {

		err = models.CategoryUpdate(p)
		if err != nil {

			_func.RespondError(context, http.StatusInternalServerError, data.Error{
				Code: fmt.Sprintf("%v", err.Error()),
				Message: "Database Error, Fail to update the category",
			})
			return
		}
		mk := make(map[string]interface{})
		mk["data"] = "category has been updated" + bson.NewObjectId().Hex()
		_func.Respond(context, http.StatusOK, mk)
	} else {
		err, categoryId := models.CategoryCreate(p)

		if err != nil {
			_func.RespondError(context, http.StatusInternalServerError, data.Error{
				Code: "502",
				Message: "Database Error, Fail to create the category",
			})
			return
		}
		mk := make(map[string]interface{})
		mk["data"] = "a category"+ categoryId + " has been successful created " + bson.NewObjectId().Hex()
		_func.Respond(context, http.StatusOK, mk)
	}
}

func CategoriesDelete(context *gin.Context) {
	id := context.PostForm("cid")
	err := models.DeleteCategory(id)
	if err.Message != "" {
		_func.RespondError(context, http.StatusInternalServerError, err)
		return
	} else {
		mk := make(map[string]interface{})
		mk["data"] = "Category Has Been Deleted! " + bson.NewObjectId().Hex()
		_func.Respond(context, http.StatusOK, mk)
	}

}
>>>>>>> nekohandserverv1/master
