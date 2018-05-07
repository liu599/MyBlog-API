package main

import (
	"os"
	"strconv"

	"nekoserver/middleware/data"
	"nekoserver/middleware/func"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	"nekoserver/router"
)


func main() {

	//Configure()
	maxIdle, _ := strconv.Atoi(os.Getenv("SERVER_DB_MAX_IDLE"))
	maxOpen, _ := strconv.Atoi(os.Getenv("SERVER_DB_MAX_OPEN"))
	source := os.Getenv("SERVER_DB_URL")

	database := data.Database{
		Driver: "mysql",
		MaxIdle: maxIdle,
		MaxOpen: maxOpen,
		Name: "bangdream",
		Source: source,
	}

	var Apps = make(map[string]data.Database)

	Apps["bangdream"] = database

	_func.AssignAppDataBaseList(Apps)

	_func.AssignDatabaseFromList([]string{"bangdream"})

	engine := gin.New()

	engine.Use(gin.Logger())

	router.AssignBackendRouter(engine)

	engine.Run(":19992")
}
