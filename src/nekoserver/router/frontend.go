package router

import (
	"github.com/gin-gonic/gin"
	"nekoserver/controller"
)

func AssignFrontendRouter(engine *gin.Engine) {
	routerGroup := engine.Group("v2/frontend")

	/*Server*/
	routerGroup.Handle("GET", "status", controller.ServerStatusGet)
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