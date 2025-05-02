package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xyperam/go-spotify-clone/models"
	"github.com/xyperam/go-spotify-clone/utils"
)

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
	var exsistingTrack models.PlaylistTrack
	if err := utils.DB.Where("spotify_id=?", input.TrackID).First(&exsistingTrack).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Track already exists in playlist"})
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

func GetPlaylistTracks(c *gin.Context) {
	var tracks []models.PlaylistTrack
	playlistIDStr := c.Param("playlistID")
	playlistID, err := strconv.Atoi(playlistIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "playlistID must be an integer"})
	}
	if err := utils.DB.First(&tracks, playlistID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "playlist not found"})
		return
	}
	if err := utils.DB.Where("playlist_id = ?", playlistID).Find(&tracks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get tracks"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"playlist_id": playlistID,
		"tracks":      tracks})
}

func GetAllPlaylist(c *gin.Context) {
	var playlists []models.Playlist
	userCtx, exist := c.Get("userID")
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userID not found in context"})
		return
	}
	userID := userCtx.(int)
	if err := utils.DB.Where("user_id = ?", userID).First(&playlists).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "playlist not found"})
		return
	}
	// Get all playlists for the user

	if err := utils.DB.Where("user_id =?", userID).Find(&playlists).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get playlists"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"playlists": playlists,
	})
}
