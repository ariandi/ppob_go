package api

import (
	"fmt"
	db "github.com/ariandi/ppob_go/db/sqlc"
	"github.com/ariandi/ppob_go/services"
	"github.com/ariandi/ppob_go/token"
	"github.com/ariandi/ppob_go/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

var userService *services.UserService

// Server serves HTTP requests for ppob services.
type Server struct {
	store      db.Store
	TokenMaker token.Maker
	Router     *gin.Engine
	config     util.Config
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		store:      store,
		TokenMaker: tokenMaker,
		config:     config,
	}
	//router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("status", validStatus)
		if err != nil {
			return nil, fmt.Errorf("cannot register validation status : %w", err)
		}
	}

	//router.POST("/users", server.createUsers)
	//router.GET("/users/:id", server.getUser)
	//router.GET("/users", server.listUsers)
	//router.PUT("/users/:id", server.updateUsers)
	//
	//server.router = router

	server.setupRouter()
	services.GetUserService(config)
	util.InitLogger()
	logrus.Println("================================================")
	logrus.Println("Server running at port %s", config.ServerAddress)
	logrus.Println("================================================")
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	router.Use(CORSMiddleware())

	router.POST("/users/login", server.loginUser)
	router.POST("/users/test-redis", server.testRedisMq)
	router.POST("/users/test-create", server.createUsersFirst)

	authRoutes := router.Group("/").Use(AuthMiddleware(server.TokenMaker))
	authRoutes.POST("/users", server.createUsers)
	authRoutes.GET("/users/:id", server.getUser)
	authRoutes.GET("/users", server.listUsers)
	authRoutes.PUT("/users/:id", server.updateUsers)

	server.Router = router
}

func (server Server) Start(address string) error {
	return server.Router.Run(address)
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
