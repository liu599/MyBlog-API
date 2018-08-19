package router

import (
	"github.com/gin-gonic/gin"

	"nekoserver/controller"
)

func AssignBackendRouter(engine *gin.Engine) {

	routerGroup := engine.Group("v2/backend")

	routerGroup.Handle("GET", "status", controller.ServerStatusGet)
	/*categories*/
	routerGroup.Handle("GET", "categories", controller.CategoriesFetch)
	/*posts*/
	routerGroup.Handle("POST", "posts", controller.PostsFetch)
	routerGroup.Handle("POST", "posts/:cid", controller.PostsFetchByCategory)
	routerGroup.Handle("POST", "post/:pid", controller.PostFetchOne)
}