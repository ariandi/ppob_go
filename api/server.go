package api

import (
	db "github.com/ariandi/ppob_go/db/sqlc"
	"github.com/gin-gonic/gin"
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

	router.POST("/users", server.createUsers)
	router.GET("/users/:id", server.getUser)
	router.GET("/users", server.listUsers)

	server.router = router
	return server
}

func (server Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}