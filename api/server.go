package api

import (
	"fmt"

	db "github.com/bacnx/simplebank/db/sqlc"
	"github.com/bacnx/simplebank/token"
	"github.com/bacnx/simplebank/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves HTTP requires for our banking service.
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer creates a new HTTP server and setup routing.
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	authGroup := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authGroup.POST("/accounts", server.createAccount)
	authGroup.GET("/accounts/:id", server.getAccount)
	authGroup.GET("/accounts", server.listAccount)
	authGroup.DELETE("/accounts/:id", server.deleteAccount)

	authGroup.POST("/transfer", server.createTransfer)

	server.router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) {
	server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
