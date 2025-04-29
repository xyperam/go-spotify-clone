package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/xyperam/go-spotify-clone/models"
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

func SearchSpotifySong(c *gin.Context) {
	query := c.Query("q")

	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query is required"})
		return
	}

	token, err := utils.GetSpotifyAccessToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get token spotify"})
		return
	}

	url := fmt.Sprintf("https://api.spotify.com/v1/search?q=%s&type=track&limit=10", url.QueryEscape(query))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
		return
	}
	req.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to execute request"})
		return
	}
	defer resp.Body.Close()
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode response"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func CreatePlaylist(c *gin.Context) {
	var inputPlaylist models.PlaylistInput

	if err := c.ShouldBindJSON(&inputPlaylist); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if inputPlaylist.PlaylistName == "" {
		inputPlaylist.PlaylistName = "My Playlist"
	}
	userCtx, exist := c.Get("userID")
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userID not found in context"})
		return
	}
	userID := userCtx.(int)

	playlist := models.Playlist{
		UserID:       userID,
		PlaylistName: inputPlaylist.PlaylistName,
	}
	result := utils.DB.Create(&playlist)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Playlist created successfully",
		"playlist": playlist})
}
