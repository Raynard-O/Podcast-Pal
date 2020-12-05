package controllers

import (
	"github.com/labstack/echo"
	"github.com/raynard2/backend/database"
	"github.com/raynard2/backend/library"
	"github.com/raynard2/backend/library/podcast"
	"github.com/raynard2/backend/library/user"
	"github.com/raynard2/backend/models"
	"log"
	"net/http"
)

func AddPodcast(c echo.Context) error {

	params := new(podcast.SearchParams)
	if err := c.Bind(params); err != nil {
		return BadRequestResponse(c, library.ErrorBinding)
	}
	DB, err := database.NewMongoConn("PodcastUsers", "users")
	if err != nil {
		log.Fatal(err)
	}

	user := user.GetUser(c)

	newPod := &models.Podcast{
		Name:     params.Name,
		PODID:       params.ID,
		ImageUrl: params.ImageUrl,
		Topics:   nil,
	}

	user.FavoritePodcast = append(user.FavoritePodcast, *newPod)

	changes := map[string][]models.Podcast{
		"FavoritePodcast": user.FavoritePodcast,
	}

	field := make(map[string]interface{})
	field["email"] = user.Email
	log.Println(changes, field)
	err = DB.UpdateOneUser(field, changes)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return c.JSON(http.StatusOK, user)
}
func FetchUserPodcast(c echo.Context) error {

	user := user.GetUser(c)

	return c.JSONPretty(http.StatusOK, user.FavoritePodcast, " ")
}
func PodcastTime(c echo.Context) error {

	params := new(podcast.PodcastParams)
	if err := c.Bind(params); err != nil {
		return BadRequestResponse(c, library.ErrorBinding)
	}
	DB, err := database.NewMongoConn("PodcastUsers", "podcast")
	if err != nil {
		log.Fatal(err)
	}
	log.Print(params)



	//user := user.GetUser(c)

	newPod := &models.Podcast{
		Name:     params.Name,
		PODID:       params.ID,
		ImageUrl: params.ImageUrl,
		Transcript: params.Transcript,
		Topics:   params.TimeStamp,
	}
	//podCheck := new(models.Podcast)
	// check if podcast exist
	_, err = DB.FindOnePodcast("podcast", "podid", params.ID)
	if err==nil {
		return MessageResponse(c, http.StatusCreated, "podcast exist in DB")
	}
	//DB.FindOne("podcast", "podid", params.ID)
	 podcastT, err := DB.SavePodcast(newPod, "podcast")
	 if err != nil {
	 	return InternalError(c, "error saving podcast timestamp: check DB")
	 }

	return DataResponse(c, http.StatusAccepted, podcastT)
}
func GetPodcast(c echo.Context) error {
	query := new(podcast.PodcastParams)
	if err := c.Bind(query); err != nil{
		return BadRequestResponse(c, library.ErrorBinding)
	}
	log.Print(query)
	Db, err := database.NewMongoConn("PodcastUsers", "podcast")
	if err != nil {
		log.Fatal(err)
	}
	podcast, err := Db.FindOnePodcast("podcast", "podid", query.ID)
	if err !=nil {
		return MessageResponse(c, http.StatusAlreadyReported, "Podcast does not exist")
	}
	return DataResponse(c, http.StatusAccepted, podcast)
}
func Timesta(c echo.Context) error {
	params := new(podcast.PodcastParam)
	if err := c.Bind(params); err != nil{
		return BadRequestResponse(c, library.ErrorBinding)
	}
	Db, err := database.NewMongoConn("PodcastUsers", "podcast")
	if err != nil {
		log.Fatal(err)
	}
	podcast, err := Db.FindOnePodcast("podcast", "podid", params.ID)
	podcast.Topics = append(podcast.Topics, params.TimeStamp)
	changes := map[string][]models.Stamp{
		"Topics": podcast.Topics,
	}
	field := make(map[string]interface{})
	field["podid"] = params.ID
	log.Println(changes, field)
	err = Db.UpdateOneUser(field, changes)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return DataResponse(c, 201, "added")
}
