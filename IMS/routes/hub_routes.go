package routes

import (
	"github.com/Trishank-Omniful/Onboarding-Task/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterHubRoutes(router *gin.RouterGroup, controller *controllers.HubController) {
	hubGroup := router.Group("/hubs")
	{
		hubGroup.GET("", controller.GetAllHubs)
		hubGroup.GET("/:id", controller.GetHubById)
		hubGroup.POST("", controller.CreateHub)
		hubGroup.PUT("/:id", controller.UpdateHub)
		hubGroup.DELETE("/:id", controller.DeleteHub)
	}
}
