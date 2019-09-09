package main

import (
	"github.com/kataras/iris"
	"github.com/xssnick/crawlyzer-auth/handlers"
	"github.com/xssnick/crawlyzer-auth/models"
	"log"
	"os"
)

func main() {
	ds, err := models.BuildStore()
	if err != nil {
		log.Println("failed to init datastore!")
		return
	}

	app := handlers.InitApp(ds)

	_ = app.Run(iris.Addr(os.Getenv("LISTEN")), iris.WithoutServerError(iris.ErrServerClosed))
}
