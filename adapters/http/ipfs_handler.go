package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) IPFSinfoHandler(c *gin.Context) {
	bootstrap, swarmKey, err := h.IS.GetConInfo()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"swarm_key": swarmKey, "bootstrap": bootstrap})
}
