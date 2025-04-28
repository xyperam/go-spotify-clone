package config

import "time"

const (
	Host     = "localhost"
	Port     = 5432
	User     = "postgres"
	Password = "123"
	DBName   = "spotify_clone"
	SSLMode  = "disable"
)

var (
	SpotifyClientID     = "b2c50d3a8ea24af592695894e4968557"
	SpotifyClientSecret = "654279b8d65e4b71bdffac139f3f1269"
	SpotifyAccessToken  string
	SpotifyTokenExpiry  time.Time
)
