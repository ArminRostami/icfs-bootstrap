// Package test _
package test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
	"testing"
	"time"

	_http "icfs_pg/adapters/http"
	db "icfs_pg/adapters/postgres"
	app "icfs_pg/application"

	. "github.com/franela/goblin"
	"github.com/gin-gonic/gin"
)

var contents = []string{
	`{
		"cid":"dsfs3mfaggasghashsgsdf6",
		"description":"a song to send instead of your projects",
		"name":"mano_yadet",
		"extension":"mp3",
		"size":25,
		"file_type":"audio"
	}`,
	`{
		"cid":"dsfssagsghejajrdafd6",
		"description":"lab report template",
		"name":"lab_report",
		"extension":"pdf",
		"size":15,
		"file_type":"document"
	}`,
}

const base = "http://127.0.0.1:8000"
const usersAPI = base + "/users"
const contentsAPI = base + "/contents"

const TypeAppJson = "application/json"

func TestE2E(t *testing.T) {
	g := Goblin(t)

	var contentIDS []string
	var pgsql *db.PGSQL

	j, err := cookiejar.New(nil)
	g.Assert(err).IsNil()
	client := &http.Client{Jar: j}

	err = os.Setenv("DEBUG", "true")
	g.Assert(err).IsNil()

	gin.SetMode(gin.ReleaseMode)

	g.Describe("Application", func() {
		g.It("should connect to database", func() {
			const conStr = "postgres://postgres:example@127.0.0.1:5432"
			var err error
			pgsql, err = db.New(conStr)
			g.Assert(err).IsNil()
		})
		g.It("should bootstrap", func() {
			go func() {
				userStore := &db.UserStore{DB: pgsql}
				contentStore := &db.ContentStore{DB: pgsql}
				contentService := &app.ContentService{CST: contentStore, UST: userStore}
				userService := &app.UserService{UST: userStore}
				handler := _http.Handler{US: userService, CS: contentService}
				handler.Serve()
			}()
		})

	})
	time.Sleep(time.Millisecond * 300)
	g.Describe("users", func() {
		g.It("should sign up", func() {
			body := []byte(`{
				"username":"testname",
				"password":"asdf",
				"email":"testmail@gmail.com"
			}`)
			resp, err := client.Post(usersAPI, TypeAppJson, bytes.NewBuffer(body))
			g.Assert(err).IsNil()
			g.Assert(resp.StatusCode).Eql(200)
		})
		g.It("should log in", func() {
			body := []byte(`{
				"username":"testname",
				"password":"asdf"
			}`)
			resp, err := client.Post(usersAPI+"/login", TypeAppJson, bytes.NewBuffer(body))
			g.Assert(err).IsNil()
			g.Assert(resp.StatusCode).Eql(200)
		})
		g.It("should update info", func() {
			payload := []byte(`{
				"email":"mailtest@yahoo.com"
			}`)
			req, err := http.NewRequest(http.MethodPut, usersAPI, bytes.NewBuffer(payload))
			g.Assert(err).IsNil()
			resp, err := client.Do(req)
			g.Assert(err).IsNil()
			g.Assert(resp.StatusCode).Eql(200)
		})
		g.It("should get info", func() {
			resp, err := client.Get(usersAPI)
			g.Assert(err).IsNil()
			g.Assert(resp.StatusCode).Eql(200)
			bytes, err := io.ReadAll(resp.Body)
			g.Assert(err).IsNil()
			var jsonObj map[string]interface{}
			err = json.Unmarshal(bytes, &jsonObj)
			g.Assert(err).IsNil()
			g.Assert(jsonObj["username"]).Eql("testname")
			g.Assert(jsonObj["email"]).Eql("mailtest@yahoo.com")
		})

	})

	g.Describe("content", func() {
		g.It("should add", func() {
			for _, c := range contents {
				cBytes := []byte(c)
				resp, err := client.Post(contentsAPI, TypeAppJson, bytes.NewBuffer(cBytes))
				g.Assert(err).IsNil()
				g.Assert(resp.StatusCode).Eql(200)
				r, err := io.ReadAll(resp.Body)
				g.Assert(err).IsNil()
				var jsonObj map[string]string
				err = json.Unmarshal(r, &jsonObj)
				g.Assert(err).IsNil()
				contentIDS = append(contentIDS, jsonObj["id"])
			}
			t.Logf("resp: %v", contentIDS)
		})
		g.It("should delete", func() {
			for _, c := range contentIDS {
				mapData := map[string]interface{}{
					"id": c,
				}
				body, err := json.Marshal(mapData)
				g.Assert(err).IsNil()
				req, err := http.NewRequest(http.MethodDelete, contentsAPI, bytes.NewBuffer(body))
				g.Assert(err).IsNil()
				resp, err := client.Do(req)
				g.Assert(err).IsNil()
				g.Assert(resp.StatusCode).Eql(200)
			}
		})
	})
	g.Describe("users", func() {
		g.It("should delete", func() {
			req, err := http.NewRequest(http.MethodDelete, usersAPI, nil)
			g.Assert(err).IsNil()
			resp, err := client.Do(req)
			g.Assert(err).IsNil()
			g.Assert(resp.StatusCode).Eql(200)

		})
	})

}
