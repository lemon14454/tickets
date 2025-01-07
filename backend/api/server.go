package api

import (
	"fmt"
	"net/http"
	"ticket/backend/cache"
	db "ticket/backend/db/sqlc"
	rbmq "ticket/backend/mq"
	"ticket/backend/token"
	"ticket/backend/util"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	config     *util.Config
	store      db.Store
	router     *gin.Engine
	tokenMaker token.TokenMaker
	mq         *rbmq.RabbitProducer
	cache      *redis.Client
}

func NewServer(config *util.Config, store db.Store) (*Server, error) {

	tokenMaker, err := token.NewJWTMaker(config.TOKEN_SYMMETRIC_KEY)
	if err != nil {
		return nil, err
	}

	mq, err := rbmq.NewRabbitProducer(config.BROKER_ADDRESS)
	if err != nil {
		return nil, err
	}
	err = mq.DeclareQueue("tickets")
	if err != nil {
		return nil, err
	}

	cache := cache.NewRedis(config.REDIS_ADDRESS)

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
		mq:         mq,
		cache:      cache,
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
	authRoutes.POST("/event", server.createEvent)
	authRoutes.POST("/ticket", server.claimTicket)
	authRoutes.POST("/order", server.createOrder)
	authRoutes.GET("/auth", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, struct{}{})
	})

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func (server *Server) Close() {
	fmt.Println("TRIGGER")
	server.mq.Close()
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
