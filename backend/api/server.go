package api

import (
	db "ticket/backend/db/sqlc"
	"ticket/backend/util"

	"github.com/gin-gonic/gin"
)

type Server struct {
	config *util.Config
	store  db.Store
	router *gin.Engine
}

func NewServer(config *util.Config, store db.Store) (*Server, error) {

	server := &Server{
		config: config,
		store:  store,
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", nil)
	router.POST("/users/login", nil)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
