package podcast

import (
	"github.com/labstack/echo"
	"github.com/raynard2/backend/models"
	"net/http"
)

type Podcast struct {
	ID                   string `json:"id"`
	Image                string `json:"image"`
	GenreIDS             []int  `json:"genre_ids"`
	Thumbnail            string `json:"thumbnail"`
	TitleOriginal        string `json:"title_original"`
	ListennotesUrl       string `json:"listennotes_url"`
	TitleHighlighted     string `json:"title_highlighted"`
	PublisherOriginal    string `json:"publisher_original"`
	PublisherHighlighted string `json:"publisher_highlighted"`
}
type Result struct {
	ID                     string  `json:"id"`
	RSS                    string  `json:"rss"`
	Link                   string  `json:"link"`
	Audio                  string  `json:"audio"`
	Image                  string  `json:"image"`
	Podcast                Podcast `json:"podcast"`
	ItuneID                int64   `json:"itunes_id"`
	Thumbnail              string  `json:"thumbnail"`
	PubDateMs              string  `json:"pub_date_ms"`
	TitleOriginal          string  `json:"title_original"`
	ListennotesUrl         string  `json:"listennotes_url"`
	AudioLengthSec         int64   `json:"audio_length_sec"`
	ExplicitContent        bool    `json:"explicit_content"`
	TitleHighlighted       string  `json:"title_highlighted"`
	DescriptionHighlighted string  `json:"description_highlighted"`
	TranscriptsHighlighted []int8  `json:"transcripts_highlighted"`
}
type PodResponse struct {
	Took       float32  `json:"took"`
	Count      int      `json:"count"`
	Total      int64    `json:"total"`
	Results    []Result `json:"results"`
	NextOffset int8     `json:"next_offset"`
}

type Response struct {
	PodcastListen   PodResponse      `json:"podcast_listen"`
	PodcastsBackend []models.Podcast `json:"podcasts_backend"`
	Success         bool             `json:"success"`
}

func PodcastResponseData(c echo.Context, PodcastListen *PodResponse, PodcastResponse []models.Podcast) error {
	response := Response{
		PodcastListen:   *PodcastListen,
		PodcastsBackend: PodcastResponse,
		Success:         true,
	}
	return c.JSONPretty(http.StatusOK, response, "")
}
