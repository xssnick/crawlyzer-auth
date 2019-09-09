package handlers

import (
	"github.com/kataras/iris"
	"github.com/xssnick/crawlyzer-auth/models"
	"net/http"
	"regexp"
)

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

func (wa *WebApp) Login(c iris.Context) {
	email := c.PostValue("email")
	pw := c.PostValue("password")

	if !emailRegex.MatchString(email) {
		c.StatusCode(http.StatusForbidden)
		c.JSON(iris.Map{
			"error": "invalid email",
		})
		return
	}

	if len(pw) < 8 {
		c.StatusCode(http.StatusForbidden)
		c.JSON(iris.Map{
			"error": "incorrect email or password",
		})
		return
	}

	ses, err := wa.Store.User.Login(email, pw)
	if err != nil {
		if err == models.ErrLoginIncorrect {
			c.StatusCode(http.StatusForbidden)
			c.JSON(iris.Map{
				"error": "incorrect email or password",
			})
		} else {
			c.StatusCode(http.StatusInternalServerError)
			c.JSON(iris.Map{
				"error": "server error",
			})
			wa.Logger.Println(err)
		}
		return
	}

	c.JSON(iris.Map{
		"success": true,
		"session":ses,
	})
}

func (wa *WebApp) Auth(c iris.Context) {
	sesid := c.PostValue("sesid")
	if sesid == "" {
		c.StatusCode(http.StatusForbidden)
		c.JSON(iris.Map{
			"success": false,
		})
		return
	}

	id, err := wa.Store.User.Auth(sesid)
	if err != nil {
		c.StatusCode(http.StatusForbidden)
		c.JSON(iris.Map{
			"success": false,
		})
		return
	}

	c.JSON(iris.Map{
		"success": true,
		"uuid":    id,
	})
}

func (wa *WebApp) RegisterNewUser(c iris.Context) {
	email := c.PostValue("email")
	pw := c.PostValue("password")

	if !emailRegex.MatchString(email) {
		c.StatusCode(http.StatusForbidden)
		c.JSON(iris.Map{
			"error": "bad email",
		})
		return
	}

	if len(pw) < 8 {
		c.StatusCode(http.StatusForbidden)
		c.JSON(iris.Map{
			"error": "bad password",
		})
		return
	}

	//TODO: check for already registered
	_, err := wa.Store.User.Create(email, pw)
	if err != nil {
		c.StatusCode(http.StatusInternalServerError)
		c.JSON(iris.Map{
			"error": "db err",
		})
		wa.Logger.Println(err)
		return
	}

	c.JSON(iris.Map{
		"code": http.StatusOK,
	})
}


func (wa *WebApp) List(c iris.Context) {
	list, err := wa.Store.User.GetAll()
	if err != nil {
		c.StatusCode(http.StatusForbidden)
		c.JSON(iris.Map{
			"success": false,
		})
		return
	}

	c.JSON(list)
}