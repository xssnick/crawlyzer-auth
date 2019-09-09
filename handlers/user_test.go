package handlers

import (
	"database/sql"
	"errors"
	"github.com/kataras/iris/httptest"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/xssnick/crawlyzer-auth/models"
	"github.com/xssnick/crawlyzer-auth/models_mock"
	"net/http"
	"testing"
)

func TestRegister(t *testing.T) {
	Convey("Register new user", t, func() {
		ds := models_mock.InitMockStore()
		ex := httptest.New(t, InitApp(ds))

		Convey("When form values are invalid", func() {
			forms := []map[string]interface{}{{
				"email":    "bop@rusggd",
				"password": "1726SU(&87h",
				"case":     "Invalid email",
				"mustbe":   "bad email",
			}, {
				"email":    "bopssa.com",
				"password": "1726SU(&87h",
				"case":     "Invalid email",
				"mustbe":   "bad email",
			}, {
				"email":    "",
				"password": "1726SU(&87h",
				"case":     "Invalid email",
				"mustbe":   "bad email",
			}, {
				"email":    "a@.com",
				"password": "1726SU(&87h",
				"case":     "Invalid email",
				"mustbe":   "bad email",
			}, {
				"email":    "@nix.com",
				"password": "322",
				"case":     "Invalid email",
				"mustbe":   "bad email",
			}, { // pw tests
				"email":    "gop@sup.com",
				"password": "123",
				"case":     "Invalid password",
				"mustbe":   "bad password",
			}, {
				"email":    "pip@sup.com",
				"password": "",
				"case":     "Invalid password",
				"mustbe":   "bad password",
			}}

			for _, variant := range forms {
				Convey("Test when: "+variant["case"].(string)+"; case - "+variant["email"].(string), func() {
					answer := ex.POST("/user/register").WithForm(variant).Expect().JSON().Object()

					So(answer.Value("error").String().Raw(), ShouldEqual, variant["mustbe"].(string))
				})
			}
		})

		Convey("When datastore errors", func() {
			forms := []map[string]interface{}{{
				"email":    "gop@sup.com",
				"password": "SuperPassword",
				"err":      sql.ErrConnDone,
			}, {
				"email":    "pop@hup.pro",
				"password": "17777223",
				"err":      errors.New("unknown"),
			}}

			for _, variant := range forms {
				ds.User.(*models_mock.MUserStore).FakeError = variant["err"].(error)

				answer := ex.POST("/user/register").WithForm(variant).Expect()

				Convey("Must be Internal error; var - "+variant["email"].(string)+" | "+variant["password"].(string), func() {
					So(answer.Raw().StatusCode, ShouldEqual, http.StatusInternalServerError)
				})
			}
		})

		Convey("When all valid", func() {
			forms := []map[string]interface{}{{
				"email":    "gop@sup.com",
				"password": "SuperPassword",
			}, {
				"email":    "pop@hup.pro",
				"password": "17777223",
			}}

			for _, variant := range forms {
				answer := ex.POST("/user/register").WithForm(variant).Expect()

				Convey("Must be OK; var - "+variant["email"].(string)+" | "+variant["password"].(string), func() {
					So(answer.Raw().StatusCode, ShouldEqual, http.StatusOK)
				})
			}
		})
	})
}

func TestLogin(t *testing.T) {
	Convey("Login user", t, func() {
		ds := models_mock.InitMockStore()
		ex := httptest.New(t, InitApp(ds))

		Convey("When form values are invalid", func() {
			forms := []map[string]interface{}{{
				"email":    "bop@rusggd",
				"password": "1726SU(&87h",
				"case":     "Invalid email",
				"mustbe":   "invalid email",
			}, {
				"email":    "bopssa.com",
				"password": "1726SU(&87h",
				"case":     "Invalid email",
				"mustbe":   "invalid email",
			}, {
				"email":    "",
				"password": "1726SU(&87h",
				"case":     "Invalid email",
				"mustbe":   "invalid email",
			}, {
				"email":    "a@.com",
				"password": "1726SU(&87h",
				"case":     "Invalid email",
				"mustbe":   "invalid email",
			}, {
				"email":    "@nix.com",
				"password": "322",
				"case":     "Invalid email",
				"mustbe":   "invalid email",
			}, { // pw tests
				"email":    "gop@sup.com",
				"password": "123",
				"case":     "Invalid password",
				"mustbe":   "incorrect email or password",
			}, {
				"email":    "pip@sup.com",
				"password": "",
				"case":     "Invalid password",
				"mustbe":   "incorrect email or password",
			}}

			for _, variant := range forms {
				Convey("Test when: "+variant["case"].(string)+"; case - "+variant["email"].(string), func() {
					answer := ex.POST("/user/login").WithForm(variant).Expect().JSON().Object()

					So(answer.Value("error").String().Raw(), ShouldEqual, variant["mustbe"].(string))
				})
			}
		})

		Convey("When datastore errors", func() {
			forms := []map[string]interface{}{{
				"email":    "gop@sup.com",
				"password": "SuperPassword",
				"err":      errors.New("unknown"),
				"mustbe":   http.StatusInternalServerError,
			}, {
				"email":    "pop@hup.pro",
				"password": "17777223",
				"err":      models.ErrLoginIncorrect,
				"mustbe":   http.StatusForbidden,
			}}

			for _, variant := range forms {
				ds.User.(*models_mock.MUserStore).FakeError = variant["err"].(error)

				answer := ex.POST("/user/login").WithForm(variant).Expect()

				Convey("Must be Internal error; var - "+variant["email"].(string)+" | "+variant["password"].(string), func() {
					So(answer.Raw().StatusCode, ShouldEqual, variant["mustbe"].(int))
				})
			}
		})

		Convey("When all valid", func() {
			forms := []map[string]interface{}{{
				"email":    "gop@sup.com",
				"password": "SuperPassword",
			}, {
				"email":    "pop@hup.pro",
				"password": "17777223",
			}}

			for _, variant := range forms {
				answer := ex.POST("/user/login").WithForm(variant).Expect()

				Convey("Must be OK; var - "+variant["email"].(string)+" | "+variant["password"].(string), func() {
					So(answer.Raw().StatusCode, ShouldEqual, http.StatusOK)
				})
			}
		})
	})
}

func TestAuth(t *testing.T) {
	Convey("Auth user", t, func() {
		ds := models_mock.InitMockStore()
		ex := httptest.New(t, InitApp(ds))

		Convey("When session is empty", func() {
			answer := ex.POST("/user/auth").WithForm(map[string]interface{}{
				"sesid": "",
			}).Expect()

			Convey("Must be auth error", func() {
				So(answer.Raw().StatusCode, ShouldEqual, http.StatusForbidden)
			})
		})

		Convey("When session is not exists", func() {
			ds.User.(*models_mock.MUserStore).FakeError = models.ErrAuthIncorrect

			answer := ex.POST("/user/auth").WithForm(map[string]interface{}{
				"sesid": "UnKnOWNsession27772",
			}).Expect()

			Convey("Must be auth error", func() {
				So(answer.Raw().StatusCode, ShouldEqual, http.StatusForbidden)
			})
		})

		Convey("When session is valid", func() {
			answer := ex.POST("/user/auth").WithForm(map[string]interface{}{
				"sesid": "UnKnOWNsession27772",
			}).Expect()

			Convey("Must be auth OK", func() {
				So(answer.Raw().StatusCode, ShouldEqual, http.StatusOK)
			})
		})
	})
}
