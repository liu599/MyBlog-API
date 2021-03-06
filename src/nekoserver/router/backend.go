package router

import (
	"github.com/gin-gonic/gin"

	"nekoserver/controller"
)

func AssignBackendRouter(engine *gin.Engine) {

	//engine.Use(auth.TokenAuthMiddleware())
	engine.Use(auth.TokenRemoteAuth())

	routerGroup := engine.Group("v2/backend")
	/*Server*/
	routerGroup.Handle("GET", "status", controller.ServerStatusGet)
	/*auth*/
	routerGroup.Handle("POST", "token.get", controller.TokenGen)
	routerGroup.Handle("POST", "token.v2.get", controller.TokenFetch)
	/*manage*/
	routerGroup.Handle("POST", "auth/post.edit", controller.PostEdit)
	routerGroup.Handle("POST", "auth/post.delete", controller.PostDelete)
	routerGroup.Handle("POST", "auth/category.edit", controller.CategoriesEdit)
	routerGroup.Handle("POST", "auth/category.delete", controller.CategoriesDelete)

}