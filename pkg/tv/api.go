package tv

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"
)
var OldMessageId int64
type Movies struct {
	Result []struct {
		ID            int    `json:"id"`
		Description   string `json:"description"`
		Votes         int    `json:"votes"`
		Author        string `json:"author"`
		Date          string `json:"date"`
		GifURL        string `json:"gifURL"`
		GifSize       int    `json:"gifSize"`
		PreviewURL    string `json:"previewURL"`
		VideoURL      string `json:"videoURL"`
		VideoPath     string `json:"videoPath"`
		VideoSize     int    `json:"videoSize"`
		Type          string `json:"type"`
		Width         string `json:"width"`
		Height        string `json:"height"`
		CommentsCount int    `json:"commentsCount"`
		FileSize      int    `json:"fileSize"`
		CanVote       bool   `json:"canVote"`
	} `json:"result"`
	TotalCount int `json:"totalCount"`
}

type Movie struct {
	Description   string `json:"description"`
	VideoPath      string `json:"videoURL"`

}
func GetMovie(log *zap.Logger) Movie {
	page := rand.Intn(100 - 1) + 1
	client := &http.Client{Timeout: 10 * time.Second}

	req , err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://developerslife.ru/top/%d?json=true", page), nil)
	if err != nil {
       log.Error("Error http NewRequest", zap.Error(err))
	}
	req.Header = http.Header{
		"Host": []string{"developerslife.ru"},
		"Content-Type": []string{"application/json"},
		"Accept": []string{"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
		"User-Agent": []string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_16_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36"},
	}
	r , err := client.Do(req)
	if err != nil {
		log.Error("Error http client", zap.Error(err))
	}
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("Error read body", zap.Error(err))
	}
	var movies Movies
	json.Unmarshal(bodyBytes, &movies)
	var movie Movie
	resultCount := len(movies.Result)
	index := rand.Intn(resultCount - 1)
	if resultCount > 0 {
		if err := DownloadFile("./tmp/video.gif", movies.Result[index].GifURL); err != nil {
			log.Error("Error DownloadFile", zap.Error(err))
		}
		movie.VideoPath = "tmp/video.gif"
		movie.Description = movies.Result[index].Description
	}
	return movie
}


func DownloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

func DeleteFile(filepath string) error {
	err := os.Remove(filepath)
	if err != nil {
		return err
	}
	return nil
}