package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/xyperam/go-spotify-clone/config"
	"github.com/xyperam/go-spotify-clone/models"
)

func GetSpotifyAccessToken() (string, error) {
	if time.Now().After(config.SpotifyTokenExpiry) {
		tokenResp, err := fetchSpotifyAccessToken()
		if err != nil {
			return "", err
		}
		config.SpotifyAccessToken = tokenResp.AccessToken
		config.SpotifyTokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	}
	return config.SpotifyAccessToken, nil
}

func fetchSpotifyAccessToken() (*models.SpotifyTokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", encodeBasicAuth(config.SpotifyClientID, config.SpotifyClientSecret))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get access token: %s", string(bodyBytes))
	}

	var tokenResp models.SpotifyTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

func encodeBasicAuth(clientID, clientSecret string) string {
	// Encode clientID and clientSecret in base64
	auth := clientID + ":" + clientSecret
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

func GetSpotifyTrackByID(trackID string, token string) (*models.SpotifyTrack, error) {
	url := fmt.Sprintf("https://api.spotify.com/v1/tracks/%s", trackID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get track details: %s", string(bodyBytes))
	}
	var track models.SpotifyTrack
	if err := json.NewDecoder(resp.Body).Decode(&track); err != nil {
		return nil, err
	}
	return &track, nil
}
