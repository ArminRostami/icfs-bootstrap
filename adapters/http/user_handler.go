package http

import (
	"icfs_mongo/domain"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const JWT = "jwt"
const ID = "id"

func (h *Handler) RegisterHandler(c *gin.Context) {
	var user domain.User
	uid := strings.Replace(uuid.New().String(), "-", "", -1)
	user.ID = uid
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

func (h *Handler) DeleteHandler(c *gin.Context) {
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

func (h *Handler) ValidateClaims(c *gin.Context) {
	id := c.GetString("id")
	c.JSON(http.StatusOK, id)
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

func (h *Handler) UpdateHandler(c *gin.Context) {
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
