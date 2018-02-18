package controller

import (
	"github.com/gin-gonic/gin"
	"nekoserver/middleware/func"
	"net/http"
)

func ServerStatusGet(context *gin.Context) {
	_func.Respond(context, http.StatusOK, gin.H{"mem":"ok"})
}