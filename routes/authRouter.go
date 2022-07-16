package routes

import (
	controller "project_login/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.LoadHTMLGlob("templates/*.html")
	incomingRoutes.GET("/login", controller.Login)
	incomingRoutes.POST("/login", controller.PostLogin)
	incomingRoutes.GET("/logout", controller.Logout)
	incomingRoutes.GET("/logoutadmin", controller.LogoutAdmin)
	incomingRoutes.GET("/signup", controller.Signup)
	incomingRoutes.POST("/signup", controller.PostSignup)
	incomingRoutes.GET("/admin", controller.Admin)
	incomingRoutes.POST("/admin", controller.PostAdmin)
	incomingRoutes.GET("/", controller.IndexHandler)
}
