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

	config := cors.DefaultConfig()

	config.AllowHeaders = []string{"User", "Origin"}

	config.AllowCredentials = true

	engine.Use(cors.New(config))

	engine.Use(static.Serve("/files/", static.LocalFile(staticRoot, true)))

	router.AssignBackendRouter(engine)

	router.AssignFrontendRouter(engine)

	engine.Run(":19992")
}
