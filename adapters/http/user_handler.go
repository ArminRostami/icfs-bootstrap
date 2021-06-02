package http

import (
	"icfs-boot/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	sessionToken = "session_token"
	userID       = "uid"
)

func (h *Handler) RegisterHandler(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, appErr := h.US.RegisterUser(&user)
	if appErr != nil {
		renderError(c, appErr)
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *Handler) DeleteUserHandler(c *gin.Context) {
	id := c.GetString(userID)

	err := h.US.DeleteUser(id)
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
	userData, sessID, err := h.US.AuthenticateUser(user.Username, user.Password)
	if err != nil {
		renderError(c, err)
		return
	}
	// TODO: determine if this is needed:
	// c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie(sessionToken, sessID, 24*3600, "/", "", false, false)
	c.JSON(http.StatusOK, gin.H{"data": userData})
}

func (h *Handler) GetUserInfo(c *gin.Context) {
	id := c.GetString(userID)
	u, err := h.US.GetUserWithID(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, u)
}

func (h *Handler) AuthorizeUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessID, err := c.Cookie(sessionToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		uid, err := h.US.ValidateAuth(sessID)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		c.Set(userID, uid)
		c.Set(sessionToken, sessID)
		c.Next()
	}
}

func (h *Handler) UserUpdateHandler(c *gin.Context) {
	id := c.GetString(userID)

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.US.UpdateUser(id, updates)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "user updated successfully"})
}

func (h *Handler) LogoutHandler(c *gin.Context) {
	sessID := c.GetString(sessionToken)

	err := h.US.Logout(sessID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "logout successful"})
}
