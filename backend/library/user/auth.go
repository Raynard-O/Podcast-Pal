package user

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/raynard2/backend/config"
	"github.com/raynard2/backend/database"
	"github.com/raynard2/backend/models"

	"log"
	"time"
)

var configuration, err = config.LoadSecrets()

var HmacSigningKey = []byte(configuration.HmacSigningKey)

func GenerateToken(user *models.User) (string, error) {
	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"user_id":   user.ID,
		"email":     user.Email,
		"issued_at": now.Unix(),
		"expire_at": now.Add(time.Hour * 72).Unix(),
	})
	return token.SignedString(HmacSigningKey)
}

func ComputeHmac256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func LoginUser(c echo.Context, user *models.User) error {
	token, _ := GenerateToken(user)
	//channel := ComputeHmac256(string(user.Model.ID), GetHmacKey)
	return LoginUserResponse(c, user, token, user.Email)
}

func GetUser(c echo.Context) *models.User {
	DB, err := database.NewMongoConn("PodcastUsers", "users")
	if err != nil {
		log.Fatal(err)
	}
	user := new(models.User)
	claims := c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)
	DB.FindOneUser("users", "email", claims["email"].(string), &user)
	//log.Println(user)
	return user
}
