package api

import (
	"fmt"
	db "github.com/ariandi/ppob_go/db/sqlc"
	"github.com/ariandi/ppob_go/middleware"
	"github.com/ariandi/ppob_go/token"
	"github.com/ariandi/ppob_go/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves HTTP requests for ppob service.
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
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users/login", server.loginUser)
	router.POST("/users/test-create", server.createUsersFirst)

	authRoutes := router.Group("/").Use(middleware.AuthMiddleware(server.TokenMaker))
	authRoutes.POST("/users", server.createUsers)
	authRoutes.GET("/users/:id", server.getUser)
	authRoutes.GET("/users", server.listUsers)
	authRoutes.PUT("/users/:id", server.updateUsers)

	server.Router = router
}

func (server Server) Start(address string) error {
	return server.Router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
