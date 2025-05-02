package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xyperam/go-spotify-clone/models"
	"github.com/xyperam/go-spotify-clone/utils"
)

func UpgradeToPremium(c *gin.Context) {
	userID := c.GetInt("userID")
	var user models.User

	if err := utils.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	user.IsPremium = true
	user.PremiumExpiresAt = time.Now().AddDate(1, 0, 0) // Set premium expiration to 1 year from now
	if err := utils.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade to premium"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully upgraded to premium"})

}
