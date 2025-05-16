//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/sidgim/finance_shared/bootstrap"
	"github.com/sidgim/finance_user/internal/server"
	"github.com/sidgim/finance_user/internal/user"
	"os"
)

// ProvideAddr lee el puerto de la env o usa ":8000" por defecto
func ProvideAddr() string {
	if a := os.Getenv("SERVER_ADDR"); a != "" {
		return a
	}
	return ":8000"
}

// 2️⃣ Módulos

var userSet = wire.NewSet(
	user.NewRepository,  // func NewRepository(*gorm.DB) *UserRepo
	user.NewService,     // func NewService(*UserRepo, *log.Logger) *Service
	user.NewUserHandler, // func NewUserHandler(*Service) *UserHandler
)

// 3️⃣ Server & Router

var serverSet = wire.NewSet(
	server.NewRouter,    // func NewRouter(*UserHandler, *CourseHandler) http.Handler
	ProvideAddr,         // func ProvideAddr() string
	bootstrap.NewServer, // func NewServer(http.Handler, addr string) *Server
)

var appSet = wire.NewSet(
	bootstrap.BaseSet,
	userSet,
	serverSet,
)

func InitializeApp() (*bootstrap.Server, error) {
	wire.Build(appSet)
	return &bootstrap.Server{}, nil
}
