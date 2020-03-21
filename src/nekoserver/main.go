package main

import (
	"os"
	"strconv"

	"nekoserver/middleware/data"
	"nekoserver/middleware/func"

<<<<<<< HEAD
=======
	"github.com/gin-contrib/cors"
>>>>>>> nekohandserverv1/master
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	"nekoserver/router"
)


func main() {

<<<<<<< HEAD
=======
	//gin.SetMode(gin.ReleaseMode)
>>>>>>> nekohandserverv1/master
	//Configure()
	maxIdle, _ := strconv.Atoi(os.Getenv("SERVER_DB_MAX_IDLE"))
	maxOpen, _ := strconv.Atoi(os.Getenv("SERVER_DB_MAX_OPEN"))
	source := os.Getenv("SERVER_DB_URL")

	database := data.Database{
		Driver: "mysql",
		MaxIdle: maxIdle,
		MaxOpen: maxOpen,
<<<<<<< HEAD
		Name: "bangdream",
=======
		Name: "nekohand",
>>>>>>> nekohandserverv1/master
		Source: source,
	}

	var Apps = make(map[string]data.Database)

<<<<<<< HEAD
	Apps["bangdream"] = database

	_func.AssignAppDataBaseList(Apps)

	_func.AssignDatabaseFromList([]string{"bangdream"})
=======
	Apps["nekohand"] = database

	_func.AssignAppDataBaseList(Apps)

	_func.AssignDatabaseFromList([]string{"nekohand"})
>>>>>>> nekohandserverv1/master

	engine := gin.New()

	engine.Use(gin.Logger())

<<<<<<< HEAD
	router.AssignBackendRouter(engine)

	engine.Run(":19992")
=======
	engine.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "OPTIONS", "PUT"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "User-Agent", "Referrer", "Host", "Token", "User", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "X-Real-Ip"},
		AllowCredentials: true,
		AllowAllOrigins:  false,
		AllowOriginFunc:  func(origin string) bool { return true },
		MaxAge:           86400,
	}))
	
	router.AssignBackendRouter(engine)

	router.AssignFrontendRouter(engine)

	engine.Run(":20479")
>>>>>>> nekohandserverv1/master
}
