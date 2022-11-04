package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/jasonLuFa/simplebank/db/sqlc"
	"github.com/jasonLuFa/simplebank/token"
	"github.com/jasonLuFa/simplebank/util"
)

type Server struct {
	config util.Config
	tokenMaker token.Maker
	store db.Store
	router *gin.Engine
}

func NewServer(config util.Config,store db.Store) (*Server,error){
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{config: config,store: store,tokenMaker: tokenMaker} 
	
	if v,ok := binding.Validator.Engine().(*validator.Validate); ok{
		v.RegisterValidation("currency",validCurrency)
		v.RegisterValidation("validAmount", validAmount)
	}
	
	server.setupRouter()
	return server,nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	router.POST("/users",server.createUser)
	router.POST("/login",server.loginUser)

	// below api need to token in header
	authoRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authoRoutes.POST("/accounts",server.createAccount)
	authoRoutes.GET("/accounts/:id",server.getAccount)
	authoRoutes.GET("/accounts",server.listAccounts)
	authoRoutes.PUT("/accounts/:id",server.updateAccount)
	authoRoutes.DELETE("/accounts/:id",server.deleteAccount)
	
	authoRoutes.POST("/transfers",server.createTransfer)
	

	server.router = router
}

func (server *Server) Start(address string) error{
	return server.router.Run(address)
}

func errorResponse(err error) gin.H{
	return gin.H{"error":err.Error()}
}