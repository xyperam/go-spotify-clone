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

	premium := protected.Group("/")
	premium.Use(middleware.RequirePremium())

	r.GET("/spotify/token", controller.GetSpotifyTokenHandler)
	r.GET("/spotify/search", controller.SearchSpotifySong)
	premium.POST("/playlist/create", controller.CreatePlaylist)
	// protected.GET("/playlist/:id/tracks", controller.AddTrackToPlaylist)
	protected.GET("/track/:trackID", controller.GetSpotifyTrackByID)
	protected.POST("/playlists/:playlistID/tracks", controller.AddTrackToPlaylist)
	r.GET("/playlists/:playlistID/tracks", controller.GetPlaylistTracks)
	protected.GET("/playlists", controller.GetAllPlaylist)
	protected.POST("/upgrade", controller.UpgradeToPremium)
	return r
}
