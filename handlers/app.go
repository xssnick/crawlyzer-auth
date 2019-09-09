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

	app.Post("/user/login", wa.Login)
	app.Post("/user/auth", wa.Auth)
	app.Post("/user/register", wa.RegisterNewUser)
	app.Get("/user/list", wa.List)
	return app
}
