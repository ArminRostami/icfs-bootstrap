package http

import (
	"icfs_cr/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) NewContentHandler(c *gin.Context) {
	var content domain.Content
	if err := c.ShouldBindJSON(&content); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	content.UploaderID = c.GetString("id")
	appErr := h.CS.RegisterContent(&content)
	if appErr != nil {
		renderError(c, appErr)
		return
	}
	c.JSON(http.StatusOK, gin.H{"cid": content.CID})
}

func (h *Handler) GetContentHandler(c *gin.Context) {
	input := struct {
		ID string `json:"id"`
	}{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	content, err := h.CS.GetContentWithID(input.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"content": content})
}

func (h *Handler) DeleteContentHandler(c *gin.Context) {
	input := struct {
		ID string `json:"id"`
	}{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid := c.GetString("id")
	err := h.CS.DeleteContent(uid, input.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "content deleted"})

}
