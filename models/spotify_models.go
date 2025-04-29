package models

import "time"

type Playlist struct {
	PlaylistID   int       `gorm:"primaryKey" json:"playlist_id"`
	UserID       int       `json:"user_id"`
	PlaylistName string    `json:"playlist_name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type PlaylistInput struct {
	PlaylistName string `json:"playlist_name" `
}
