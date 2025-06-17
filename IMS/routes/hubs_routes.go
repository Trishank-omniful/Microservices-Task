package routes

import (
	"github.com/Trishank-Omniful/Onboarding-Task/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterHubRoutes(router *gin.RouterGroup) {
	hubGroup := router.Group("/hubs")
	{
		hubGroup.GET("", controllers.GetAllHubs)
		hubGroup.GET("/:id", controllers.GetHubById)
		hubGroup.POST("", controllers.CreateHub)
		// hubGroup.PUT("/:id", controllers.UpdateHub)
		// hubGroup.DELETE("/:id", controllers.DeleteHub)
	}
}
