package main

import (
	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"lifresh/api"
	"time"
)

func init() {

}

func main() {

	router := gin.Default()

	router.Use(cors.New(
		cors.Config{
			AllowAllOrigins: true,
			AllowedMethods:  []string{"POST, GET"},
			MaxAge:          12 * time.Hour,
		}))

	//router.Use(cors.Middleware(cors.Config{
	//	Origins:         "*",
	//	Methods:         "GET, PUT, POST, DELETE",
	//	RequestHeaders:  "Origin, Authorization, Content-Type",
	//	ExposedHeaders:  "",
	//	MaxAge:          50 * 1000,
	//	Credentials:     true,
	//	ValidateHeaders: false,
	//}))
	//router.Use(cors.Default())
	router.Use(static.Serve("/", static.LocalFile("../html", true)))
	// router.LoadHTMLGlob("templates/*")
	//router.GET("/", hello)
	router.POST("/api/:name", api.ApiCall)

	router.Run(":8000")
}
