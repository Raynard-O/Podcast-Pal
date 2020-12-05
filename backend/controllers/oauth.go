package controllers

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"github.com/raynard2/backend/config"
	"github.com/raynard2/backend/database"
	user2 "github.com/raynard2/backend/library/user"
	"github.com/raynard2/backend/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"strings"
	"time"
)

type Credentials struct {
	ID     string `json:"id"`
	Secret string `json:"secret"`
}

// User is a retrieved and authentiacted user.
type User struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Gender        string `json:"gender"`
}

var state string
var conf *oauth2.Config
var configuration, err = config.LoadSecrets()

func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

var (
	// TODO: randomize it
	oauthStateString = randToken()
)
//const scopes = [
//	'https://www.googleapis.com/auth/userinfo.profile',
//	'https://www.googleapis.com/auth/userinfo.email',
//];
var local = 1;
func HandleGoogleLogin(c echo.Context) error {
	var _ string
	if local == 1 {
		_ = "http://127.0.0.1:5001/google/auth"
	}else {
		_ = "http://podcast.ca-central-1.elasticbeanstalk.com/google/auth"
	}

	conf = &oauth2.Config{
		ClientID:     configuration.GoogleCientID,
		ClientSecret: configuration.GoogleSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://127.0.0.1:5001/google/auth",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile", // You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
		},
	}
	url := conf.AuthCodeURL(oauthStateString)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func Goredirect(c echo.Context) error {
	var htmlIndex = `<html>
<body>
	<a href="/googlelogin">Google Log In</a>
</body>
</html>`
	return c.HTML(200, htmlIndex)
}

func getuserinfo(state string, code string) ([]byte, *oauth2.Token, error) {
	if state != oauthStateString {
		return nil, nil, fmt.Errorf("invalid oauth state")
	}
	token, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}
	return contents, token, nil
}

func HandleGoogleCallback(c echo.Context) error {
	_, token, err := getuserinfo(c.FormValue("state"), c.FormValue("code"))
	if err != nil {
		fmt.Println(err.Error())
		return c.Redirect(http.StatusTemporaryRedirect, "/google")
	}

	googleClient := conf.Client(oauth2.NoContext, token)
	resp, err := googleClient.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return BadRequestResponse(c, err.Error())
	}
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)
	log.Println("Resp body: ", string(data))
	Guser := new(User)
	_ = json.Unmarshal(data, &Guser)

	DB, err := database.NewMongoConn("PodcastUsers", "users")
	if err != nil {
		log.Fatal(err)
	}
	if !Guser.EmailVerified {
		return BadRequestResponse(c, "google Auth 'not verified' ")
	}

	user := &models.User{
		ID:        primitive.NewObjectID(),
		Sub:       Guser.Sub,
		Username:  Guser.GivenName,
		FullName:  Guser.Name,
		Email:     Guser.Email,
		CreatedAt: time.Now(),
		Active:    true,
		Method:    "google",
	}
	//check if user details already exist
	if Guser.Sub != "" {
		if _, err = DB.FindOne("users", "sub", Guser.Sub); err != nil {
			_, err = DB.Save(user, "users")
			if err != nil {
				log.Fatal(err.Error())
			}
			log.Print("google user created ")
			//return BadRequestResponse(c, "OAUTH Error")
		}
	} else {
		log.Print("google user exist ")
	}
	c.Set("user_g", Guser.Sub)

	var url string
	if local == 1 {
		url = "http://042e6caaca17.ngrok.io/guser" + "?user_id=" + "Guser.Sub"
	}else {
		url = "http://podpal.ca-central-1.elasticbeanstalk.com/guser" + "?user_id=" + Guser.Sub
	}

	log.Print(url)
	//go ReadCookie(c)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func GiveToken(c echo.Context) error {
	query := c.QueryParam("user_id")
	log.Println(query)
	user := new(models.User)
	DB, err := database.NewMongoConn("PodcastUsers", "users")
	if err != nil {
		log.Fatal(err)
	}
	if user, err = DB.FindOne("users", "sub", query); err != nil {
		_, err = DB.Save(user, "users")
		if err != nil {
			log.Fatal(err.Error())
		}
		log.Print("google user created")
		//return BadRequestResponse(c, "OAUTH Error")
	}
	log.Print(user)
	//token, _ := user2.GenerateToken(user)
	//user.Token = token
	//return DataResponse(c, http.StatusContinue, user)
	return user2.LoginUser(c, user)
}

func ReadCookie(c echo.Context) error {
	query := c.QueryParam("user_id")
	log.Print(query)
	user := new(models.User)
	DB, err := database.NewMongoConn("PodcastUsers", "users")
	if err != nil {
		log.Fatal(err)
	}
	if user, err = DB.FindOne("users", "sub", query); err != nil {
		_, err = DB.Save(user, "users")
		if err != nil {
			log.Fatal(err.Error())
		}
		log.Print("google user created ")
		//return BadRequestResponse(c, "OAUTH Error")
	}
	log.Print(user)
	token, _ := user2.GenerateToken(user)
	url := "http://podpal.ca-central-1.elasticbeanstalk.com/account"
	cookie := new(http.Cookie)

	cookie.Name = "user_google"
	cookie.Value = user.Email + `|` + token
	cookie.Expires = time.Now().Add(30 * time.Second)
	cookie.HttpOnly = true
	c.SetCookie(cookie)
	return c.Redirect(http.StatusMovedPermanently, url)
}
func Read(c echo.Context) error {
	cookie2, err := c.Cookie("user_google")
	if err != nil {
		return BadRequestResponse(c, err.Error())
	}
	xs := strings.Split(cookie2.Value, "|")
	fmt.Println("HERE'S THE SLICE", xs)
	email := xs[0]
	codeRcvd := xs[1]
	codeCheck := getCode(email)
	fmt.Println(codeRcvd)
	fmt.Println(codeCheck)
	return c.JSON(200, cookie2.Value)
}

func getCode(data string) string {
	h := hmac.New(sha256.New, []byte("ourkey"))
	io.WriteString(h, data)
	return fmt.Sprintf("%x", h.Sum(nil))
}
