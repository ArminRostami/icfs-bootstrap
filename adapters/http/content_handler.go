package http

import (
	"icfs_pg/domain"
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
	id, appErr := h.CS.RegisterContent(&content)
	if appErr != nil {
		renderError(c, appErr)
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *Handler) GetContentHandler(c *gin.Context) {
	id := c.Query("id")
	uid := c.GetString("id")
	content, err := h.CS.GetContentWithID(uid, id)
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

func (h *Handler) ContentUpdateHandler(c *gin.Context) {
	id := c.GetString(ID)

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.CS.UpdateContent(id, updates)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "content updated successfully"})
}

func (h *Handler) SearchHandler(c *gin.Context) {
	var search map[string]string
	if err := c.ShouldBindJSON(&search); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	results, err := h.CS.SearchContent(search)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, results)
}

func (h *Handler) RateContentHandler(c *gin.Context) {
	input := struct {
		Rating float32 `json:"rating"`
		CID    string  `json:"content_id"`
	}{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid := c.GetString("id")
	err := h.CS.RateContent(input.Rating, uid, input.CID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "rating submitted."})
}

func (h *Handler) TextSearchHandler(c *gin.Context) {
	input := struct {
		Term string `json:"term"`
	}{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	results, err := h.CS.TextSearch(input.Term)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"results": results})

}
