// Package http includes the http handler and routes
package http

import (
	app "icfs_pg/application"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type Handler struct {
	ge *gin.Engine
	US *app.UserService
	CS *app.ContentService
}

func (h *Handler) Serve() error {
	h.ge = gin.Default()
	h.ge.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:4200"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Set-Cookie", "Origin", "Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
		ExposeHeaders:    []string{"Set-Cookie"},
	}))
	h.SetupRoutes()
	err := h.ge.Run(":8000")
	return errors.Wrap(err, "failed to start gin engine")
}

func (h *Handler) SetupRoutes() {
	h.ge.POST("/users", h.RegisterHandler)
	h.ge.GET("/users", h.AuthorizeJWT(), h.GetUserInfo)
	h.ge.PUT("/users", h.AuthorizeJWT(), h.UserUpdateHandler)
	h.ge.DELETE("/users", h.AuthorizeJWT(), h.DeleteUserHandler)

	h.ge.POST("/users/login", h.LoginHandler)

	h.ge.POST("/contents", h.AuthorizeJWT(), h.NewContentHandler)
	h.ge.GET("/contents", h.AuthorizeJWT(), h.GetContentHandler)
	h.ge.PUT("/contents", h.AuthorizeJWT(), h.ContentUpdateHandler)
	h.ge.DELETE("/contents", h.AuthorizeJWT(), h.DeleteContentHandler)

	h.ge.POST("/contents/rate", h.AuthorizeJWT(), h.RateContentHandler)

	h.ge.GET("/contents/all", h.GetAllContentsHandler)
	h.ge.POST("/contents/search", h.TextSearchHandler)
}

func renderError(c *gin.Context, appErr *app.Error) {
	c.JSON(appErr.Status, gin.H{"error": appErr.Err.Error()})
}
