package http

import (
	"icfs_cr/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

const JWT = "jwt"
const ID = "id"

func (h *Handler) RegisterHandler(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, appErr := h.USV.RegisterUser(&user)
	if appErr != nil {
		renderError(c, appErr)
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *Handler) DeleteUserHandler(c *gin.Context) {
	id := c.GetString("id")

	err := h.USV.DeleteUser(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "user deleted successfully"})
}

func (h *Handler) LoginHandler(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tok, err := h.USV.AuthenticateUser(user.Username, user.Password)
	if err != nil {
		renderError(c, err)
		return
	}
	c.SetCookie(JWT, tok, 24*3600, "/", "", true, false)
	c.JSON(http.StatusOK, gin.H{"username": user.Username})
}

func (h *Handler) GetUserInfo(c *gin.Context) {
	id := c.GetString("id")
	u, err := h.USV.GetUserWithID(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, u)
}

func (h *Handler) AuthorizeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwt, err := c.Cookie(JWT)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		claims, err := h.USV.ValidateAuth(jwt)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		c.Set(ID, claims.ID)
		c.Next()
	}
}

func (h *Handler) UserUpdateHandler(c *gin.Context) {
	id := c.GetString(ID)

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.USV.UpdateUser(id, updates)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "user updated successfully"})
}

func (h *Handler) SearchHandler(c *gin.Context) {
	var search map[string]string
	if err := c.ShouldBindJSON(&search); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if _, exists := search["term"]; !exists {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": `key "term" does not exist`})
		return
	}
	results, err := h.USV.SearchInBio(search["term"])
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, results)
}
