package handlers

import (
	"github.com/kataras/iris"
	"github.com/xssnick/crawlyzer-auth/models"
	"log"
	"os"
)

type WebApp struct {
	Store  *models.DataStore
	Logger *log.Logger
}

var app *iris.Application

func InitApp(store *models.DataStore) *iris.Application {
	app = iris.Default()

	wa := &WebApp{
		Store:  store,
		Logger: log.New(os.Stdout, "[handler]", log.LstdFlags|log.Lshortfile),
	}

	app.Post("/login", wa.Login)
	app.Post("/auth", wa.Auth)
	app.Post("/register", wa.RegisterNewUser)
	app.Get("/list", wa.List)
	app.Get("/node", wa.Node)
	app.OnErrorCode(404,func(c iris.Context) {
		c.JSON(c.Request().URL.String())
	})
	return app
}
