package router

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/raynard2/backend/config"
	"github.com/raynard2/backend/controllers"
	"net/http"
)

var configuration, err = config.LoadSecrets()
var HmacSigningKey = []byte(configuration.HmacSigningKey)

func New() *echo.Echo {
	e := echo.New()
	_ = e.Group("/ee")
	e.Static("/static", "assets")
	//library.InitClient()
	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS512",
		SigningKey:    HmacSigningKey,
		Skipper: func(c echo.Context) bool {
			if c.Path() == "/" || c.Path() == "/g" || c.Path() == "/login" || c.Path() == "/signup" || c.Path() == "/games" || c.Path() == "/google" || c.Path() == "/googlelogin" || c.Path() == "/google/auth" || c.Path() == "/google/guser" || c.Path() == "/guser" {
				return true
			}
			return false
		},
	}))
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	//templates

	//// servers other static files
	//staticBox := rice.MustFindBox("../static")
	//staticFileServer := http.StripPrefix("/static/", http.FileServer(staticBox.HTTPBox()))
	//e.GET("/static/*", echo.WrapHandler(staticFileServer))
	//
	////Set Renderer
	//e.Renderer = gorice.New(rice.MustFindBox("../views"))
	//
	//e.GET("/games", func(c echo.Context) error {
	//	//render only file, must full name with extension
	//	return c.Render(http.StatusOK, "newrpsTL.html", echo.Map{"title": "Page file title!!"})
	//})


	//entry point/homepage
	e.GET("/", func(context echo.Context) error {
		return context.JSONPretty(http.StatusOK, "Home", " ")
	})
	//authentication
	e.POST("/signup", controllers.CreateUser)
	e.POST("/login", controllers.Login)
	e.GET("/p", controllers.Podcast)
	//google OAUTH api
	e.GET("/g", controllers.ReadCookie)
	e.GET("/google", controllers.Goredirect)
	e.GET("/googlelogin", controllers.HandleGoogleLogin)
	e.GET("/google/auth", controllers.HandleGoogleCallback)
	e.GET("/guser", controllers.GiveToken)
	// podcast api
	e.POST("/favorite", controllers.AddPodcast)
	e.GET("/getfavorite", controllers.FetchUserPodcast)
	e.POST("/savetopics", controllers.PodcastTime)
	e.GET("/gettopics", controllers.GetPodcast)
	e.POST("/addtopic", controllers.Timesta)
	//GET users API
	e.GET("/getuser", controllers.GetUsersInfo)

	return e
}
