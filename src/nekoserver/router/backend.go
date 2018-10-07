package router

import (
	"github.com/gin-gonic/gin"
	"nekoserver/middleware/auth"

	"nekoserver/controller"
)

func AssignBackendRouter(engine *gin.Engine) {

	engine.Use(auth.TokenAuthMiddleware())

	routerGroup := engine.Group("v2/backend")
	/*Server*/
	routerGroup.Handle("GET", "status", controller.ServerStatusGet)
	/*Auth*/
	//routerGroup.Handle("POST", "token.get", controller.TokenGenerator)
	/*categories*/
	routerGroup.Handle("GET", "categories", controller.CategoriesFetch)
	/*posts*/
	routerGroup.Handle("POST", "posts", controller.PostsFetch)
	routerGroup.Handle("POST", "posts/:cid", controller.PostsFetchByCategory)
	routerGroup.Handle("POST", "post/:pid", controller.PostFetchOne)
	routerGroup.Handle("GET", "posts-chronology", controller.PostsChornology)
	/*comments*/
	routerGroup.Handle("POST", "comments/:pid", controller.CommentsFetch)
	routerGroup.Handle("POST", "c2a5cc3b070", controller.CommentCreation)
	/*static file*/
}