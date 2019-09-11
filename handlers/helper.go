package handlers

import (
	"github.com/kataras/iris"
)

func ThrowError(c iris.Context, code int, text string) {
	c.StatusCode(code)
	c.JSON(iris.Map{
		"error": text,
	})
}
