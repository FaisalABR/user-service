package controllers

import (
	controllers "user-service/controllers/user"
	"user-service/services"
)

type Registry struct {
	service services.IServiceRegistery
}

type IControllerRegistry interface {
	GetUserController() controllers.IUserController
}

func NewControllerRegistry(service services.IServiceRegistery) IControllerRegistry {
	return &Registry{service: service}
}

func (u *Registry) GetUserController() controllers.IUserController {
	return controllers.NewUserController(u.service)
}
