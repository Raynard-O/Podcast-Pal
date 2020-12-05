package user

import (
	"github.com/labstack/echo"
	"github.com/raynard2/backend/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type UserResponse struct {
	ID       primitive.ObjectID
	Email    string `json:"email"`
	Username string `json:"username"`
	Channel  string `json:"channel"`
	Active   bool   `json:"active"`
}
type Response struct {
	User    UserResponse `json:"user"`
	Token   string       `json:"token"`
	Success bool         `json:"success"`
}

func UserResponseData(c echo.Context, user *models.User, token string, channel string) error {
	response := Response{
		User: UserResponse{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
			Channel:  channel,
			Active:   user.Active,
		},
		Token:   token,
		Success: true,
	}

	return c.JSONPretty(http.StatusOK, response, "")
}

func LoginUserResponse(c echo.Context, user *models.User, token string, channel string) error {
	response := Response{
		User: UserResponse{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
			Channel:  channel,
			Active:   user.Active,
		},
		Token:   token,
		Success: true,
	}

	return c.JSONPretty(http.StatusOK, response, "")
}
