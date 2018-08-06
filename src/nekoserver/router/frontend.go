package router

import (
	"nekoserver/controller"

	"github.com/gin-gonic/gin"
)

func AssignFrontendRouter(engine *gin.Engine) {
	routerGroup := engine.Group("v2/frontend")

	/*categories*/
	routerGroup.Handle("GET", "categories", controller.CategoriesFetch)

}