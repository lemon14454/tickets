package api

import (
	"ticket/backend/cache"
	db "ticket/backend/db/sqlc"
	rbmq "ticket/backend/mq"
	"ticket/backend/token"
	"ticket/backend/util"

	"github.com/gin-gonic/gin"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

var WaitExchange = "wait_exchange"

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

	err = mq.DeclareExchange(TicketExchange, rbmq.DirectRouting)
	if err != nil {
		return nil, err
	}

	err = mq.DeclareExchange(WaitExchange, rbmq.DirectRouting)
	if err != nil {
		return nil, err
	}

	err = mq.DeclareQueue(TicketQueue, amqp091.Table{})
	if err != nil {
		return nil, err
	}

	err = mq.QueueBind(TicketQueue, TicketExchange, TicketRoutingKey)
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

	router.GET("/event", server.listEvent)
	router.GET("/event/:id", server.listEventZone)

	// authRoutes := router.Group("/").Use(AuthMiddleware(server.tokenMaker), RateLimitMiddleware(server.cache))
	authRoutes := router.Group("/").Use(AuthMiddleware(server.tokenMaker)) 
	authRoutes.POST("/event", server.createEvent)
	authRoutes.POST("/ticket", server.claimTicket)

	authRoutes.POST("/order", server.createOrder)
	authRoutes.GET("/order", server.listOrder)
	authRoutes.GET("/order/:id", server.orderDetail)


	rateLimitRoutes := router.Group("/limit").Use(AuthMiddleware(server.tokenMaker), RateLimitMiddleware(server.cache))

	rateLimitRoutes.POST("/ticket", server.claimTicket)
	rateLimitRoutes.POST("/order", server.createOrder)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func (server *Server) Close() {
	server.mq.Close()
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
