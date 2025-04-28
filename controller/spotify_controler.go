package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xyperam/go-spotify-clone/utils"
)

func GetSpotifyTokenHandler(c *gin.Context) {
	// Implement logic to get Spotify token
	token, err := utils.GetSpotifyAccessToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"access_token": token})
}
