package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"pet-spotlight/io"
	"strings"
)

type playerResponse struct {
	StreamingData streamingData `json:"streamingData"`
}

type streamingData struct {
	Formats         []format `json:"formats"`
	AdaptiveFormats []format `json:"adaptiveFormats"`
}

type format struct {
	ITag itag   `json:"itag"`
	URL  string `json:"url"`
}

type itag int

const (
	mp4360p itag = 18
	mp4720p itag = 22
)

// DownloadVideo downloads the Youtube video.
func DownloadVideo(youtubeURL string, path string, fileName string) error {
	downloadURL, err := getDownloadURL(youtubeURL)
	if err != nil {
		return err
	}
	decodedDownloadURL, err := url.QueryUnescape(downloadURL)
	if err != nil {
		return err
	}
	resp, err := http.Get(decodedDownloadURL)
	if err != nil {
		return err
	}
	defer io.CloseResource(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to download video")
	}
	filePath := fmt.Sprintf("%s/%s", path, fileName)
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filePath, err)
	}
	defer io.CloseResource(f)
	if err = io.CopyToFile(resp.Body, f); err != nil {
		return err
	}
	return err
}

func getDownloadURL(youtubeURL string) (string, error) {
	u := fmt.Sprintf("http://youtube.com/get_video_info?video_id=%s", strings.Split(youtubeURL, "?v=")[1][:11])
	resp, err := http.Get(u)
	if err != nil {
		return "", err
	}
	defer io.CloseResource(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to get video info")
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	params := strings.Split(string(b), "&")
	var encodedPlayerResponse string
	for _, param := range params {
		if strings.HasPrefix(param, "player_response") {
			encodedPlayerResponse = param
			break
		}
	}
	playerResponseJSON, err := url.QueryUnescape(encodedPlayerResponse[16:])
	if err != nil {
		return "", err
	}
	var pr playerResponse
	if err = json.Unmarshal([]byte(playerResponseJSON), &pr); err != nil {
		return "", err
	}
	var downloadURL string
	for _, format := range pr.StreamingData.Formats {
		if format.ITag == mp4720p {
			downloadURL = format.URL
			break
		}
		if len(downloadURL) == 0 && format.ITag == mp4360p {
			downloadURL = format.URL
		}
	}
	return downloadURL, nil
}
