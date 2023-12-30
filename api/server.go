package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/kanishkmittal55/simplebank/db/sqlc"
	"github.com/kanishkmittal55/simplebank/db/util"
	"github.com/kanishkmittal55/simplebank/token"
)

// Server serves HTTP request for our banking service.
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer creates a new HTTP server instance and setup routing.
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token : %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	// Custom Validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	// Allowing CORS Request ( TODO: Further Tweaking Required )
	router.Use(CORSMiddleware())

	// Landing Page
	router.GET("/landing/:key", server.getLandingPage)

	// Routes for Users ( No Login Required , because we want everyone to freely create an account on our server
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	router.POST("/token/renew_access", server.renewAccessToken)

	authRoute := router.Group("/").Use(authMiddleware(server.tokenMaker))

	// Add routes to router
	authRoute.POST("/accounts", server.createAccount)
	authRoute.GET("/accounts/:id", server.getAccountById)
	authRoute.GET("/accounts", server.ListAccount)
	authRoute.PUT("/accounts/:id", server.UpdateAccountById)
	authRoute.DELETE("/accounts/:id", server.deleteAccount)

	// Routes for Transfer
	authRoute.POST("/transfers", server.createTransfer)

	// router.GET("/users/:usr", server.getUserByUsername)

	server.router = router

}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
