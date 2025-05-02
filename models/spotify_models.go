package models

import "time"

type Playlist struct {
	PlaylistID   int             `gorm:"primaryKey" json:"playlist_id"`
	UserID       int             `json:"user_id"`
	PlaylistName string          `json:"playlist_name"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	Tracks       []PlaylistTrack `gorm:"foreignKey:PlaylistID" json:"tracks,omitempty"`
}

type PlaylistInput struct {
	PlaylistName string `json:"playlist_name" `
}

type PlaylistTrack struct {
	PlaylistID int    `json:"-"`
	Title      string `json:"title"`
	Artist     string `json:"artist"`
	Album      string `json:"album"`
	SpotifyID  string `json:"spotify_id"`
	PreviewURL string `json:"preview_url"`
}

type SpotifyTrack struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Album struct {
		Name string `json:"name"`
	} `json:"album"`
	Artists []struct {
		Name string `json:"name"`
	} `json:"artists"`
	PreviewURL string `json:"preview_url"`
}

type InputTrackToPlaylist struct {
	PlaylistID int    `json:"playlist_id"`
	TrackID    string `json:"track_id"`
}
