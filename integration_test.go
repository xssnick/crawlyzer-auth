// +build integration

package main

import (
	uuid "github.com/iris-contrib/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/xssnick/crawlyzer-auth/models"
	"testing"
)

func bootstrap(name string, f func(ds *models.DataStore), t *testing.T) {
	Convey(name, t, func() {
		ds,err := models.BuildStore()
		if err != nil {
			t.Fatal("failed to init datastore")
			return
		}

		f(ds)

		models.MigrateDown(ds.Postgres.DB)
	})
}

func TestCreate(t *testing.T) {
	bootstrap("User creation", func(ds *models.DataStore) {
		Convey("When OK", func() {
			uid,err:=ds.User.Create("kis@pips.com", "231vre423")

			So(err, ShouldEqual, nil)
			So(uid, ShouldNotEqual, uuid.Nil)
		})

		Convey("When duplicate", func() {
			uid,err:=ds.User.Create("kis@pips.com", "231vre423")

			So(err, ShouldEqual, nil)
			So(uid, ShouldNotEqual, uuid.Nil)

			_,err = ds.User.Create("kis@pips.com", "3423423")
			So(err, ShouldEqual, models.ErrAlreadyCreated)
		})
	}, t)
}

func TestLogin(t *testing.T) {
	bootstrap("User login", func(ds *models.DataStore) {
		Convey("When user not exists", func() {
			_, err := ds.User.Login("kis@pips.com", "123456789")
			So(err, ShouldEqual, models.ErrLoginIncorrect)
		})

		Convey("When password incorrect", func() {
			ds.User.Create("kis@pips.com", "123456789")

			_, err := ds.User.Login("kis@pips.com", "BadPassword")
			So(err, ShouldEqual, models.ErrLoginIncorrect)
		})

		Convey("When all correct", func() {
			ds.User.Create("kis@pips.com", "123456789")

			ses, err := ds.User.Login("kis@pips.com", "123456789")
			So(err, ShouldEqual, nil)
			So(ses, ShouldNotBeBlank)
		})
	}, t)
}

func TestGetAll(t *testing.T) {
	bootstrap("Users list", func(ds *models.DataStore) {
		Convey("When 2 users", func() {
			ds.User.Create("kis@pips.com", "7564756fg")
			ds.User.Create("poo@six.biz", "12346453FFF")

			all, err := ds.User.GetAll()
			So(err, ShouldEqual, nil)
			So(len(all), ShouldEqual, 2)
		})

		Convey("When empty", func() {
			all, err := ds.User.GetAll()
			So(err, ShouldEqual, nil)
			So(len(all), ShouldEqual, 0)
		})
	}, t)
}

func TestAuth(t *testing.T) {
	bootstrap("User auth", func(ds *models.DataStore) {
		Convey("When session not exists", func() {
			_,err:=ds.User.Auth("unknown-ses-id")
			So(err, ShouldEqual, models.ErrAuthIncorrect)
		})

		Convey("When session exists", func() {
			cuid,_:=ds.User.Create("kis@pips.com", "7564756fg")
			sesid,_:=ds.User.Login("kis@pips.com", "7564756fg")

			uid,err:=ds.User.Auth(sesid)
			So(err, ShouldEqual, nil)
			So(uid, ShouldEqual, cuid)
		})
	}, t)
}
