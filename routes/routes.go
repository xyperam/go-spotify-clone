package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/xyperam/go-spotify-clone/controller"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()

	r.POST("/register", controller.RegisterUser)
	r.POST("/login", controller.LoginUser)
	return r
}
