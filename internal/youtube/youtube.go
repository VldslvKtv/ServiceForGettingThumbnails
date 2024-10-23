package youtube

import (
	"encoding/json"
	"fmt"
	"getthumbnails/internal/config"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

const pathToApiData = "./config/api.yaml"

type Thumbnail struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type Thumbnails struct {
	Default  Thumbnail `json:"default"`
	Medium   Thumbnail `json:"medium"`
	High     Thumbnail `json:"high"`
	Standard Thumbnail `json:"standard"`
	Maxres   Thumbnail `json:"maxres"`
}

type Snippet struct {
	PublishedAt  string     `json:"publishedAt"`
	ChannelID    string     `json:"channelId"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	Thumbnails   Thumbnails `json:"thumbnails"`
	ChannelTitle string     `json:"channelTitle"`
	Tags         []string   `json:"tags"`
}

type Item struct {
	ID      string  `json:"id"`
	Snippet Snippet `json:"snippet"`
}

type Response struct {
	Items []Item `json:"items"`
}

func getVideoID(videoURL string) (string, error) {
	const op = "youtube.getVideoID"
	u, err := url.Parse(videoURL)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	query := u.Query()
	return query.Get("v"), nil
}

func getThumbnailURL(videoID string, apiURL string, apiKey string) (string, error) {
	const op = "youtube.getThumbnailURL"
	url := fmt.Sprintf("%s?key=%s&part=snippet&id=%s", apiURL, apiKey, videoID)
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("no items found for video ID: %s", videoID)
	}

	if len(response.Items) == 0 {
		return "", fmt.Errorf("no items found for video ID: %s", videoID)
	}

	fmt.Println(response.Items[0].Snippet.Thumbnails.High.URL)

	return response.Items[0].Snippet.Thumbnails.High.URL, nil
}

func downloadThumbnail(thumbnailURL, filename string) error {
	const op = "youtube.downloadThumbnail"

	thumbnailsDir := "thumbnails"
	err := os.MkdirAll(thumbnailsDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	filepath := filepath.Join(thumbnailsDir, filename)

	resp, err := http.Get(thumbnailURL)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer file.Close()

	_, err = file.ReadFrom(resp.Body)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func GetThumbnail(videoUrl string) (string, error) {
	apiData := config.ApiKeyAndUrl(pathToApiData)
	videoID, err := getVideoID(videoUrl)
	if err != nil {
		return "", err
	}
	thumbnailURL, err := getThumbnailURL(videoID, apiData.APIURL, apiData.APIKey)
	if err != nil {
		return "", err
	}

	if err := downloadThumbnail(thumbnailURL, fmt.Sprintf("%s.jpg", videoID)); err != nil {
		return "", err
	}

	return thumbnailURL, nil

}
