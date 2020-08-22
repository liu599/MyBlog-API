package controller

import (
	"net/http"

	"gopkg.in/mgo.v2/bson"
	"nekoserver/middleware/data"
	"nekoserver/middleware/func"

	"github.com/gin-gonic/gin"
)


func ServerStatusGet(context *gin.Context) {

	m := data.Post{}
	pg := data.Pager{}


	m.Author = "eddie32"
	m.PTitle = "Post Title"
	m.Body = bson.NewObjectId().Hex()
	m.Category = "adjfkasjflkajkdla"

	pg.PageSize = 15
	pg.PageNum = 1
	pg.TotalNumber = 1500

	mk := make(map[string]interface{})
	mk["data"] = m
	mk["pager"] = pg

	_func.Respond(context, http.StatusOK, gin.H{"mem":"ok"}, mk)
}