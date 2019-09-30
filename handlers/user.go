package handlers

import (
	"github.com/kataras/iris"
	"github.com/xssnick/crawlyzer-auth/models"
	"net/http"
	"os"
	"regexp"
)

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

func (wa *WebApp) Login(c iris.Context) {
	email := c.PostValue("email")
	pw := c.PostValue("password")

	if !emailRegex.MatchString(email) {
		ThrowError(c, http.StatusForbidden, "invalid email")
		return
	}

	if len(pw) < 8 {
		ThrowError(c, http.StatusForbidden, "incorrect email or password")
		return
	}

	ses, err := wa.Store.User.Login(email, pw)
	if err != nil {
		if err == models.ErrLoginIncorrect {
			ThrowError(c, http.StatusForbidden, "invalid email")
		} else {
			wa.Logger.Println(err)
			ThrowError(c, http.StatusInternalServerError, "server error")
		}
		return
	}

	c.JSON(iris.Map{
		"session": ses,
	})
}

func (wa *WebApp) Auth(c iris.Context) {
	sesid := c.PostValue("sesid")
	if sesid == "" {
		ThrowError(c, http.StatusForbidden, "incorrect session")
		return
	}

	id, err := wa.Store.User.Auth(sesid)
	if err != nil {
		ThrowError(c, http.StatusForbidden, "incorrect session")
		return
	}

	c.JSON(iris.Map{
		"uuid": id,
	})
}

func (wa *WebApp) RegisterNewUser(c iris.Context) {
	email := c.PostValue("email")
	pw := c.PostValue("password")

	if !emailRegex.MatchString(email) {
		ThrowError(c, http.StatusForbidden, "bad email")
		return
	}

	if len(pw) < 8 {
		ThrowError(c, http.StatusForbidden, "bad password")
		return
	}

	//TODO: check for already registered
	_, err := wa.Store.User.Create(email, pw)
	if err != nil {
		wa.Logger.Println(err)
		ThrowError(c, http.StatusInternalServerError, "server error")
		return
	}

	c.JSON(iris.Map{
		"success": true,
	})
}

func (wa *WebApp) List(c iris.Context) {
	list, err := wa.Store.User.GetAll()
	if err != nil {
		wa.Logger.Println(err)
		ThrowError(c, http.StatusInternalServerError, "server error")
		return
	}

	c.JSON(list)
}

func (wa *WebApp) Node(c iris.Context) {
	c.JSON(os.Getenv("NODEID"))
}
