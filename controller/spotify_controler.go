package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

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

func AddTrackToPlaylist(c *gin.Context) {
	playlistIDStr := c.Param("playlistID")
	playlistID, err := strconv.Atoi(playlistIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "playlistID must be an integer"})
		return
	}

	var input struct {
		TrackID string `json:"track_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	track, err := utils.FetchSpotifyTrackByID(input.TrackID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get track"})
		return
	}
	playlistTrack := models.PlaylistTrack{
		PlaylistID: playlistID,
		Title:      track.Name,
		Artist:     track.Artists[0].Name,
		Album:      track.Album.Name,
		SpotifyID:  track.ID,
		PreviewURL: track.PreviewURL,
	}
	result := utils.DB.Create(&playlistTrack)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// Kembalikan response
	c.JSON(http.StatusOK, gin.H{"message": "Track added to playlist successfully", "track": playlistTrack})
}

func GetSpotifyTrackByID(c *gin.Context) {
	// Ambil trackID dari URL parameter
	trackID := c.Param("trackID")
	if trackID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "trackID is required"})
		return
	}

	// Debugging: Log trackID
	fmt.Println("trackID:", trackID)

	// Ambil token Spotify
	token, err := utils.GetSpotifyAccessToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get token spotify"})
		return
	}

	// Bangun URL untuk request ke Spotify API
	url := fmt.Sprintf("https://api.spotify.com/v1/tracks/%s", trackID)
	// Debugging: Log URL request
	fmt.Println("Request URL:", url)

	// Buat request ke Spotify API
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
		return
	}

	// Set header Authorization
	req.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to execute request"})
		return
	}
	defer resp.Body.Close()

	// Debugging: Log response status code
	fmt.Println("Response Status Code:", resp.StatusCode)

	// Cek apakah status code OK
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("failed to get track: %s, Status: %d", string(bodyBytes), resp.StatusCode),
		})
		return
	}

	// Decode response dari Spotify API
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode response"})
		return
	}

	// Return response yang berhasil
	c.JSON(http.StatusOK, result)
}
