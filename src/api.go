package src

import (
	"TRAFforward/src/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Api(r *gin.Engine, db *gorm.DB) {

	// r.POST("/users", func(c *gin.Context) {
	// 	var user models.User
	// 	if err := c.ShouldBindJSON(&user); err != nil {
	// 		c.JSON(400, gin.H{"error": err.Error()})
	// 		return
	// 	}
	// 	db.Create(&user)
	// 	c.JSON(200, user)
	// })
	r.POST("/transferred", func(ctx *gin.Context) {

		go RunTransferred(0, "127.0.0.1:30003", "127.0.0.1:57890")

		var users []models.User
		db.Find(&users)
		ctx.JSON(200, users)
	})

	r.GET("/users", func(c *gin.Context) {
		var users []models.User
		db.Find(&users)
		c.JSON(200, users)
	})
}
