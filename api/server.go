package api

import (
	db "github.com/ariandi/ppob_go/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves HTTP requests for ppob service.
type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("status", validStatus)
		if err != nil {
			return nil
		}
	}

	router.POST("/users", server.createUsers)
	router.GET("/users/:id", server.getUser)
	router.GET("/users", server.listUsers)
	router.PUT("/users/:id", server.updateUsers)

	server.router = router
	return server
}

func (server Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
