package svc

import (
	"api/internal/config"
	"api/internal/middleware"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config    config.Config
	AdminAuth rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		AdminAuth: middleware.NewAdminAuthMiddleware().Handle,
	}
}
