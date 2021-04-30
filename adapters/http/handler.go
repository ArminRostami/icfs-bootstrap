// Package http includes the http handler and routes
package http

import (
	"icfs_pg/adapters/ipfs"
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
	IS *ipfs.IpfsService
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

func renderError(c *gin.Context, appErr *app.Error) {
	c.JSON(appErr.Status, gin.H{"error": appErr.Err.Error()})
}
