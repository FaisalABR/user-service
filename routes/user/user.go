package routes

import (
	"user-service/controllers"
	"user-service/middlewares"

	"github.com/gin-gonic/gin"
)

type UserRoute struct {
	controllers controllers.IControllerRegistry
	group       *gin.RouterGroup
}

type IUserRoute interface {
	Run()
}

func NewUserRoute(controllers controllers.IControllerRegistry, group *gin.RouterGroup) IUserRoute {
	return &UserRoute{
		controllers: controllers,
		group:       group,
	}
}

func (u *UserRoute) Run() {
	group := u.group.Group("/auth")
	group.GET("/user", middlewares.Authenticate(), u.controllers.GetUserController().GetUserLogin)
	group.GET("/:uuid", middlewares.Authenticate(), u.controllers.GetUserController().GetUserByUUID)
	group.POST("/login", u.controllers.GetUserController().Login)
	group.POST("/register", u.controllers.GetUserController().Register)
	group.PUT("/:uuid", middlewares.Authenticate(), u.controllers.GetUserController().Update)
}
