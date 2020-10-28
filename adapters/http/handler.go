// Package http includes the http handler and routes
package http

import (
	app "icfs_mongo/application"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type Handler struct {
	ge  *gin.Engine
	USV *app.UserService
}

func (h *Handler) Serve() error {
	h.ge = gin.Default()
	h.SetupRoutes()
	err := h.ge.Run()
	return errors.Wrap(err, "failed to start gin engine")
}

func (h *Handler) SetupRoutes() {
	h.ge.POST("/register", h.RegisterHandler)
	h.ge.POST("/login", h.LoginHandler)
	h.ge.GET("/validate", h.AuthorizeJWT(), h.ValidateClaims)
}

func renderError(c *gin.Context, appErr *app.Error) {
	c.JSON(appErr.Status, gin.H{"error": appErr.Err.Error()})
}
