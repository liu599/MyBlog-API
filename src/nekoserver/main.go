package main

import (
	"os"
	"strconv"

	"nekoserver/middleware/data"
	"nekoserver/middleware/func"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	"nekoserver/router"
)


func main() {

	gin.SetMode(gin.ReleaseMode)
	//Configure()
	maxIdle, _ := strconv.Atoi(os.Getenv("SERVER_DB_MAX_IDLE"))
	maxOpen, _ := strconv.Atoi(os.Getenv("SERVER_DB_MAX_OPEN"))
	source := os.Getenv("SERVER_DB_URL")
	staticRoot := os.Getenv("SERVER_STATIC_ROOT")

	database := data.Database{
		Driver: "mysql",
		MaxIdle: maxIdle,
		MaxOpen: maxOpen,
		Name: "nekohand",
		Source: source,
	}

	var Apps = make(map[string]data.Database)

	Apps["nekohand"] = database

	_func.AssignAppDataBaseList(Apps)

	_func.AssignDatabaseFromList([]string{"nekohand"})

	engine := gin.New()

	engine.Use(gin.Logger())

	engine.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "OPTIONS", "PUT"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "User-Agent", "Referrer", "Host", "Token", "User"},
		ExposeHeaders:    []string{"Content-Length", "X-Real-Ip"},
		AllowCredentials: true,
		AllowAllOrigins:  false,
		AllowOriginFunc:  func(origin string) bool { return true },
		MaxAge:           86400,
	}))

	engine.Use(static.Serve("/files/", static.LocalFile(staticRoot, true)))

	router.AssignBackendRouter(engine)

	router.AssignFrontendRouter(engine)

	engine.Run(":19992")
}
