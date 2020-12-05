package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"github.com/raynard2/backend/database"
	"github.com/raynard2/backend/library"
	"github.com/raynard2/backend/library/podcast"
	"github.com/raynard2/backend/library/user"
	userLib "github.com/raynard2/backend/library/user"
	"github.com/raynard2/backend/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	listenApi = os.Getenv("LISTEN_NOTE_API")
	url       = os.Getenv("LISTEN_NOTE_URL")
)

func CreateUser(c echo.Context) error {
	DB, err := database.NewMongoConn("PodcastUsers", "users")

	if err != nil {
		log.Fatal(err)
	}
	params := new(userLib.CreateUserParams)
	// bind params
	err = c.Bind(params)
	if err != nil {

		return BadRequestResponse(c, err.Error())
	}

	if params.Email == "" || params.Password == "" || params.ConfirmPassword == "" || params.FullName == "" || params.Username == "" {
		return BadRequestResponse(c, library.EMPTY)
	}
	//check if user details already exist
	if params.Email != "" {
		//if  err = DB.FindOneUser("users", "email", params.Email, nil); err == nil {
		//	return BadRequestResponse(c, library.EmailTaken)
		//}

		if _, err = DB.FindOne("users", "email", params.Email); err == nil {
			return BadRequestResponse(c, library.EmailTaken)
		}
		log.Println(err)
	}

	if params.Username != "" {
		//if  err = DB.FindOneUser("users", "username", params.Username, nil); err == nil {
		//	return BadRequestResponse(c, library.EmailTaken)
		//}
		if _, err = DB.FindOne("users", "email", params.Email); err == nil {
			return BadRequestResponse(c, library.UsernameTaken)
		}
	}
	hash, err := library.GenerateHash(params.Password)
	user := models.User{
		ID:        primitive.NewObjectID(),
		Username:  params.Username,
		FullName:  params.FullName,
		Email:     params.Email,
		Password:  hash,
		CreatedAt: time.Now(),
		Method:    "local",
		Active:    true,
	}
	users, err := DB.Save(user, "users")
	if err != nil {
		log.Fatal(err.Error())
	}
	return userLib.UserResponseData(c, &users, "", user.Email)
}

func Login(c echo.Context) error {
	DB, err := database.NewMongoConn("PodcastUsers", "users")
	if err != nil {
		log.Fatal(err)
	}
	params := new(userLib.LoginParams)
	err = c.Bind(params)
	if err != nil {
		return BadRequestResponse(c, err.Error())
	}

	userData := new(models.User)
	if params.Email != "" {
		fmt.Print(params.Email)
		DB.FindOneUser("users", "username", params.Username, &userData)
		//user, err := DB.FindOne("users", "email", params.Email)
		if err != nil {
			return BadRequestResponse(c, err.Error())
		}
		//userData = user
	} else if params.Username != "" {
		DB.FindOneUser("users", "username", params.Username, &userData)
		//user, err := DB.FindOne("users", "username", params.Email)
		if err != nil {
			return BadRequestResponse(c, err.Error())
		}
		//userData = user
	}
	passBool := library.CompareHashWithPassword(userData.Password, params.Password)
	log.Print(passBool)
	if !passBool {
		return BadRequestResponse(c, library.WrongPassword)
	}

	return userLib.LoginUser(c, userData)
}

func Podcast(c echo.Context) error {
	client := http.Client{
		//Timeout: time.Duration(20 * time.Second),
	}
	// fetch listen note api key from env
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Fatal(err)
	}
	// set header
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("X-ListenAPI-Key", listenApi)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// decode body response into params
	params := new(podcast.PodResponse)
	json.NewDecoder(resp.Body).Decode(&params)
	if err != nil {
		log.Fatal(err)
	}

	log.Print(params.Results[1].Audio)

	var FFullurl = params.Results[5].Audio

	request, err := client.Get(FFullurl)
	defer request.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	b, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Panic(err)
	}
	l := bytes.NewBuffer(b)
	// log size of file before putting in aws s3 bucket
	fmt.Printf("Just Downloaded a file %s with size %d\n", "file", int64(len(b)))
	file_name := []string{params.Results[3].TitleOriginal, "mp3"}

	go library.AwsUpload(l, int64(len(b)), strings.Join(file_name, "."))
	//bytes.NewBuffer()
	var podmongo []models.Podcast

	return podcast.PodcastResponseData(c, params, podmongo)
}

//
//func P(c echo.Context) error {
//	params:= new(podcast.PodResponse)
//	var FFullurl string = params.Results[5].Audio
//	content, size,  err  := library.Getl(url, listen_api, FFullurl)
//	if err != nil {
//		log.Panic(err)
//	}
//
//	fmt.Printf("Just Downloaded a file %s with size %d\n", "file", size)
//	file_name := []string{params.Results[1].TitleOriginal, "mp3"}
//
//	go library.AwsUpload(content , size, strings.Join(file_name, "."))
//	//bytes.NewBuffer()
//	var podmongo []models.Podcast
//
//
//	return podcast.PodcastResponseData(c, params, podmongo)
//}
//

func GetUsersInfo(c echo.Context) error {
	users := user.GetUser(c)
	return c.JSONPretty(http.StatusAccepted, users, " ")
}
