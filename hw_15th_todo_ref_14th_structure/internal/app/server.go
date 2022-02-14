package app

import (
	"todolist/internal/api"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router          *gin.Engine
	todolistService api.TodolistService
	loginServive    api.LoginService
}

func NewServer(r *gin.Engine, ts api.TodolistService, ls api.LoginService) *Server {
	return &Server{
		router:          r,
		todolistService: ts,
		loginServive:    ls,
	}
}

func (s *Server) Run() error {
	r := s.Routers()
	err := endless.ListenAndServe(":9010", r)
	if err != nil {
		return err
	}
	return nil
}
