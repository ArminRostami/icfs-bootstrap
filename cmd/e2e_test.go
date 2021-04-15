package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
	"testing"
	"time"

	_http "icfs_pg/adapters/http"
	db "icfs_pg/adapters/postgres"
	app "icfs_pg/application"
	"icfs_pg/domain"

	. "github.com/franela/goblin"
	"github.com/gin-gonic/gin"
)

var mockContent1 = []map[string]interface{}{
	{
		"cid":         "dsfs3mfaggasghashsgsdf6",
		"description": "a song to send instead of your projects",
		"name":        "mano_yadet",
		"extension":   "mp3",
		"size":        25,
		"file_type":   "audio",
	},
	{
		"cid":         "dsfssagsghejajrdafd6",
		"description": "lab report template",
		"name":        "lab_report",
		"extension":   "pdf",
		"size":        15,
		"file_type":   "document",
	},
	{
		"cid":         "dsfsajdffjadkdkakajrdafd6",
		"description": "internship report template",
		"name":        "internship_report",
		"extension":   "docx",
		"size":        35,
		"file_type":   "document",
	},
	{
		"cid":         "dsfscderhajeafjrdafd6",
		"description": "a famous farsi font",
		"name":        "bnazanin",
		"extension":   "ttf",
		"size":        5,
		"file_type":   "font",
	},
}
var mockContent2 = []map[string]interface{}{
	{
		"cid":         "dsfs3mfgsagsahasahsgsdf6",
		"description": "a james bond movie",
		"name":        "quantom_of_solace",
		"extension":   "mkv",
		"size":        1225,
		"file_type":   "video",
	},
	{
		"cid":         "dsfssagsahasgjsagjsjkdafd6",
		"description": "blade trilogy collection",
		"name":        "blade_collection",
		"extension":   "zip",
		"size":        3215,
		"file_type":   "archive",
	},
	{
		"cid":         "dsfssagsahamthsngnsddafd6",
		"description": "the original windows xp wallpaper",
		"name":        "win_xp_wallpaper",
		"extension":   "jpeg",
		"size":        2,
		"file_type":   "image",
	},
	{
		"cid":         "dsfssagnsgsndgnbynrdafd6",
		"description": "intended for crawlers only :)",
		"name":        "robots",
		"extension":   "txt",
		"size":        2,
		"file_type":   "text",
	},
}

var users = []string{
	`{
		"username":"testname",
		"password":"asdf",
		"email":"testmail@gmail.com"
	}`,
	`{
		"username":"mrtester",
		"password":"asdf",
		"email":"mrtester@gmail.com"
	}`,
}

const base = "http://127.0.0.1:8000"
const usersAPI = base + "/users"
const contentsAPI = base + "/contents"

const TypeAppJson = "application/json"

func TestE2E(t *testing.T) {
	g := Goblin(t)

	var contentIDS []string
	var contentIDS2 []string
	var pgsql *db.PGSQL

	jar1, err := cookiejar.New(nil)
	g.Assert(err).IsNil()
	client1 := &http.Client{Jar: jar1}

	jar2, err := cookiejar.New(nil)
	g.Assert(err).IsNil()
	client2 := &http.Client{Jar: jar2}

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
	g.Describe("user1", func() {
		g.It("should sign up", func() {
			body := []byte(users[0])
			resp, err := client1.Post(usersAPI, TypeAppJson, bytes.NewBuffer(body))
			g.Assert(err).IsNil()
			g.Assert(resp.StatusCode).Eql(200)
		})
		g.It("should log in", func() {
			body := []byte(`{
				"username":"testname",
				"password":"asdf"
			}`)
			resp, err := client1.Post(usersAPI+"/login", TypeAppJson, bytes.NewBuffer(body))
			g.Assert(err).IsNil()
			g.Assert(resp.StatusCode).Eql(200)
		})
		g.It("should update info", func() {
			payload := []byte(`{
				"email":"mailtest@yahoo.com"
			}`)
			req, err := http.NewRequest(http.MethodPut, usersAPI, bytes.NewBuffer(payload))
			g.Assert(err).IsNil()
			resp, err := client1.Do(req)
			g.Assert(err).IsNil()
			g.Assert(resp.StatusCode).Eql(200)
		})
		g.It("should add content", func() {
			for _, c := range mockContent1 {
				cBytes, err := json.Marshal(c)
				g.Assert(err).IsNil()
				resp, err := client1.Post(contentsAPI, TypeAppJson, bytes.NewBuffer(cBytes))
				g.Assert(err).IsNil()
				g.Assert(resp.StatusCode).Eql(200)
				r, err := io.ReadAll(resp.Body)
				g.Assert(err).IsNil()
				var jsonObj map[string]string
				err = json.Unmarshal(r, &jsonObj)
				g.Assert(err).IsNil()
				contentIDS = append(contentIDS, jsonObj["id"])
			}
		})
		g.It("should get info", func() {
			resp, err := client1.Get(usersAPI)
			g.Assert(err).IsNil()
			g.Assert(resp.StatusCode).Eql(200)
			bytes, err := io.ReadAll(resp.Body)
			g.Assert(err).IsNil()
			var jsonObj map[string]interface{}
			err = json.Unmarshal(bytes, &jsonObj)
			g.Assert(err).IsNil()
			credit := 0
			for _, c := range mockContent1 {
				credit += c["size"].(int)
			}
			g.Assert(jsonObj["credit"]).Eql(float64(credit))
			g.Assert(jsonObj["username"]).Eql("testname")
			g.Assert(jsonObj["email"]).Eql("mailtest@yahoo.com")
		})

	})

	g.Describe("user2", func() {
		g.It("should sign up", func() {
			body := []byte(users[1])
			resp, err := client2.Post(usersAPI, TypeAppJson, bytes.NewBuffer(body))
			g.Assert(err).IsNil()
			g.Assert(resp.StatusCode).Eql(200)
		})
		g.It("should log in", func() {
			body := []byte(`{
				"username":"mrtester",
				"password":"asdf"
			}`)
			resp, err := client2.Post(usersAPI+"/login", TypeAppJson, bytes.NewBuffer(body))
			g.Assert(err).IsNil()
			g.Assert(resp.StatusCode).Eql(200)
		})
		g.It("should add content", func() {
			for _, c := range mockContent2 {
				cBytes, err := json.Marshal(c)
				g.Assert(err).IsNil()
				resp, err := client2.Post(contentsAPI, TypeAppJson, bytes.NewBuffer(cBytes))
				g.Assert(err).IsNil()
				g.Assert(resp.StatusCode).Eql(200)
				r, err := io.ReadAll(resp.Body)
				g.Assert(err).IsNil()
				var jsonObj map[string]string
				err = json.Unmarshal(r, &jsonObj)
				g.Assert(err).IsNil()
				contentIDS2 = append(contentIDS2, jsonObj["id"])
			}
		})
		g.It("should get content", func() {
			resp, err := client2.Get(contentsAPI + "?id=" + contentIDS[0])
			g.Assert(err).IsNil()
			g.Assert(resp.StatusCode).Eql(200)
		})
		g.It("should rate", func() {
			body := []byte(fmt.Sprintf(`
			{
				"rating":4.6,
				"content_id":"%s"
			}`, contentIDS[0]))
			resp, err := client2.Post(contentsAPI+"/rate", TypeAppJson, bytes.NewBuffer(body))
			g.Assert(err).IsNil()
			bts, err := io.ReadAll(resp.Body)
			t.Logf("rate resp: %v", string(bts))
			g.Assert(err).IsNil()
			g.Assert(resp.StatusCode).Eql(200)
		})
		g.It("should comment on content", func() {
			body := []byte(fmt.Sprintf(`{
				"id":"%s",
				"comment":"terrible stuff"
			}`, contentIDS[0]))
			resp, err := client2.Post(contentsAPI+"/comment", TypeAppJson, bytes.NewBuffer(body))
			g.Assert(err).IsNil()
			g.Assert(resp.StatusCode).Eql(200)
		})
		g.It("should get info", func() {
			resp, err := client2.Get(usersAPI)
			g.Assert(err).IsNil()
			g.Assert(resp.StatusCode).Eql(200)
			bytes, err := io.ReadAll(resp.Body)
			g.Assert(err).IsNil()
			var jsonObj map[string]interface{}
			err = json.Unmarshal(bytes, &jsonObj)
			g.Assert(err).IsNil()
			credit := 0
			for _, c := range mockContent2 {
				credit += c["size"].(int)
			}
			credit = credit - mockContent1[0]["size"].(int)
			g.Assert(jsonObj["credit"]).Eql(float64(credit))
		})
		g.It("should read comments", func() {
			resp, err := client2.Get(contentsAPI + "/comment?id=" + contentIDS[0])
			g.Assert(err).IsNil()
			g.Assert(resp.StatusCode).Eql(200)
			bts, err := io.ReadAll(resp.Body)
			g.Assert(err).IsNil()
			var jsonObj []domain.Comment
			err = json.Unmarshal(bts, &jsonObj)
			g.Assert(err).IsNil()
			g.Assert(len(jsonObj) <= 0).IsFalse()
			g.Assert(jsonObj[0].CText).Eql("terrible stuff")
		})
		g.Xit("should delete contents", func() {
			for _, c := range contentIDS2 {
				mapData := map[string]interface{}{
					"id": c,
				}
				body, err := json.Marshal(mapData)
				g.Assert(err).IsNil()
				req, err := http.NewRequest(http.MethodDelete, contentsAPI, bytes.NewBuffer(body))
				g.Assert(err).IsNil()
				resp, err := client2.Do(req)
				g.Assert(err).IsNil()
				g.Assert(resp.StatusCode).Eql(200)
			}
		})
		g.Xit("should delete account", func() {
			req, err := http.NewRequest(http.MethodDelete, usersAPI, nil)
			g.Assert(err).IsNil()
			resp, err := client2.Do(req)
			g.Assert(err).IsNil()
			g.Assert(resp.StatusCode).Eql(200)

		})

	})
	g.Describe("user1", func() {
		g.Xit("should delete contents", func() {
			for _, c := range contentIDS {
				mapData := map[string]interface{}{
					"id": c,
				}
				body, err := json.Marshal(mapData)
				g.Assert(err).IsNil()
				req, err := http.NewRequest(http.MethodDelete, contentsAPI, bytes.NewBuffer(body))
				g.Assert(err).IsNil()
				resp, err := client1.Do(req)
				g.Assert(err).IsNil()
				g.Assert(resp.StatusCode).Eql(200)
			}
		})
		g.Xit("should delete account", func() {
			req, err := http.NewRequest(http.MethodDelete, usersAPI, nil)
			g.Assert(err).IsNil()
			resp, err := client1.Do(req)
			g.Assert(err).IsNil()
			g.Assert(resp.StatusCode).Eql(200)

		})
	})

}
