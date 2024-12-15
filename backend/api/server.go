package api

import (
	"net/http"
	db "ticket/backend/db/sqlc"
	"ticket/backend/token"
	"ticket/backend/util"

	"github.com/gin-gonic/gin"
)

type Server struct {
	config     *util.Config
	store      db.Store
	router     *gin.Engine
	tokenMaker token.TokenMaker
}

func NewServer(config *util.Config, store db.Store) (*Server, error) {

	tokenMaker, err := token.NewJWTMaker(config.TOKEN_SYMMETRIC_KEY)
	if err != nil {
		return nil, err
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	router.POST("/tokens/renew", server.renewToken)

	authRoutes := router.Group("/").Use(AuthMiddleware(server.tokenMaker))
	authRoutes.GET("/auth", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, struct{}{})
	})

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
