package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/xyperam/go-spotify-clone/controller"
	"github.com/xyperam/go-spotify-clone/middleware"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()

	r.POST("/register", controller.RegisterUser)
	r.POST("/login", controller.LoginUser)
	protected := r.Group("/")
	protected.Use(middleware.JWTMiddleWare())
	r.GET("/spotify/token", controller.GetSpotifyTokenHandler)
	r.GET("/spotify/search", controller.SearchSpotifySong)
	protected.POST("/playlist/create", controller.CreatePlaylist)
	// protected.GET("/playlist/:id/tracks", controller.AddTrackToPlaylist)
	protected.GET("/track/:trackID", controller.GetSpotifyTrackByID)

	return r
}
