package main

import (

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	"nekoserver/router"
)


func main() {

	engine := gin.New()

	engine.Use(gin.Logger())

	router.AssignBackendRouter(engine)

	engine.Run(":19992")
}
