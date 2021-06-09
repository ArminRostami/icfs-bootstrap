package http

import (
	"icfs-boot/domain"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) NewContentHandler(c *gin.Context) {
	var content domain.Content
	if err := c.ShouldBindJSON(&content); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	content.UploaderID = c.GetString(userID)
	id, appErr := h.CS.RegisterContent(&content)
	if appErr != nil {
		renderError(c, appErr)
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *Handler) GetContentHandler(c *gin.Context) {
	id := c.Query("id")
	uid := c.GetString(userID)
	content, err := h.CS.GetContentWithID(uid, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println(content)
	c.JSON(http.StatusOK, gin.H{"content": content})
}

func (h *Handler) DeleteContentHandler(c *gin.Context) {
	content_id := c.Query("id")
	uid := c.GetString(userID)
	err := h.CS.DeleteContent(uid, content_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "content deleted"})

}
func (h *Handler) DeleteDownloadHandler(c *gin.Context) {
	content_id := c.Query("id")
	uid := c.GetString(userID)
	err := h.CS.DeleteDownload(uid, content_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "content deleted"})
}

func (h *Handler) ContentUpdateHandler(c *gin.Context) {
	id := c.GetString(userID)

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

func (h *Handler) ReviewContentHandler(c *gin.Context) {
	input := struct {
		CID     string  `json:"content_id"`
		Rating  float32 `json:"rating"`
		Comment string  `json:"comment"`
	}{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid := c.GetString(userID)
	err := h.CS.AddReview(uid, input.CID, input.Comment, input.Rating)
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

func (h *Handler) GetAllContentsHandler(c *gin.Context) {
	results, err := h.CS.GetAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"results": results})
}

func (h *Handler) GetUserUploadsHandler(c *gin.Context) {
	uid := c.GetString(userID)
	results, err := h.CS.GetUserUploads(uid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"results": results})
}

func (h *Handler) GetUserDownloadsHandler(c *gin.Context) {
	uid := c.GetString(userID)
	results, err := h.CS.GetUserDownloads(uid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"results": results})
}

func (h *Handler) GetCommentsHandler(c *gin.Context) {
	id := c.Query("id")
	comments, err := h.CS.GetComments(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, comments)
}
