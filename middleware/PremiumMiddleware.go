package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xyperam/go-spotify-clone/models"
	"github.com/xyperam/go-spotify-clone/utils"
)

func RequirePremium() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt("userID")
		var user models.User
		if err := utils.DB.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			c.Abort()
			return
		}
		if !user.IsPremium || user.PremiumExpiresAt.Before(time.Now()) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Premium access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}
