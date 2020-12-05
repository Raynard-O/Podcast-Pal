package podcast

import "github.com/raynard2/backend/models"

type SearchParams struct {
	Name     string `json:"name"`
	ID       string `json:"id"`
	ImageUrl string `json:"image_url"`
}




type PodcastParams struct {
	Name     string `json:"name"`
	ID       string `json:"id"`
	ImageUrl string `json:"image_url"`
	Transcript string 	`json:"transcript"`
	TimeStamp []models.Stamp	`json:"time_stamp"`
}


type PodcastParam struct {
	Name     string `json:"name"`
	ID       string `json:"id"`
	ImageUrl string `json:"image_url"`
	Transcript string 	`json:"transcript"`
	TimeStamp models.Stamp	`json:"time_stamp"`
}


//
//type timestamp struct {
//	time 	string
//	topics	string
//}
