package http

import (
	"fmt"
	"icfs_mongo/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

const AuthHeader = "Authorization"
const BearerSchema = "Bearer"

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

func (h *Handler) DeleteHandler(c *gin.Context) {
	id := c.GetString("id")
	// if !exists {
	// 	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no claims"})
	// 	return
	// }
	// clmap, ok := claims.(app.CustomClaims)
	// if !ok {
	// 	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "cannot cast claims to map"})
	// 	return
	// }
	// fmt.Println(clmap.ID)

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
	c.Header(AuthHeader, fmt.Sprintf("%s %s", BearerSchema, tok))
}

func (h *Handler) ValidateClaims(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no claims"})
	}
	c.JSON(http.StatusOK, claims)
}

func (h *Handler) AuthorizeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AuthHeader)
		tokenStr := authHeader[len(BearerSchema)+1:]
		claims, err := h.USV.ValidateAuth(tokenStr)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		c.Set("id", claims.ID)
		c.Next()
	}
}
