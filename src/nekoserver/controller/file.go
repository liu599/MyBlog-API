package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"nekoserver/middleware/func"
	"nekoserver/models"
)

func FileListFetch(c *gin.Context) {
	err, filelist := models.FetchFileList()
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Error filelist %s", err.Error()))
		return
	}
	list := make(map[string]interface{})
	list["data"] = filelist
	_func.Respond(c, http.StatusOK, list)
}
