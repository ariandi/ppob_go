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

var userService services.UserInterface
var roleService services.RoleInterface
var categoryService services.CategoryInterface
var partnerService services.PartnerInterface
var providerService services.ProviderInterface
var productService services.ProductInterface
var transactionService services.TransactionInterface
var sellingService *services.SellingService
var roleUserService services.RoleUserInterface
var mediaService services.MediaInterface
var paymentService services.PaymentInterface

// Server serves HTTP requests for ppob services.
type Server struct {
	//store      db.Store
	tokenMaker token.Maker
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
		//store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("status", validStatus)
		if err != nil {
			return nil, fmt.Errorf("cannot register validation status : %w", err)
		}

		err = v.RegisterValidation("paymentType", validPaymentType)
		if err != nil {
			return nil, fmt.Errorf("cannot register validation status : %w", err)
		}
	}

	server.setupRouter(tokenMaker)
	userService = services.GetUserService(config, store, tokenMaker)
	roleService = services.GetRoleService(store)
	categoryService = services.GetCategoryService(store)
	partnerService = services.GetPartnerService(store)
	providerService = services.GetProviderService(store)
	productService = services.GetProductService(store)
	transactionService = services.GetTransactionService(store)
	sellingService = services.GetSellingService(store)
	roleUserService = services.GetRoleUserService(store)
	mediaService = services.GetMediaService(store)
	paymentService = services.GetPaymentService(store)
	util.InitLogger()
	logrus.Println("================================================")
	logrus.Printf("Server running at port %s", config.ServerAddress)
	logrus.Println("================================================")
	return server, nil
}

func (server *Server) setupRouter(tokenMaker token.Maker) {
	router := gin.Default()
	router.Use(CORSMiddleware())

	router.POST("/users/login", server.loginUser)
	router.POST("/users/test-redis", server.testRedisMq)
	router.POST("/users/test-create", server.createUsersFirst)
	router.GET("/transactions/payment-export", server.paymentExport)
	router.POST("/digi-webhook", server.digiWebHook)

	authRoutes := router.Group("/").Use(AuthMiddleware(tokenMaker))
	authRoutes.POST("/users", server.createUsers)
	authRoutes.GET("/users/:id", server.getUser)
	authRoutes.GET("/users", server.listUsers)
	authRoutes.PUT("/users/:id", server.updateUsers)
	authRoutes.DELETE("/users/:id", server.softDeleteUser)

	authRoutes.POST("/roles", server.createRole)
	authRoutes.GET("/roles/:id", server.getRole)
	authRoutes.GET("/roles", server.listRole)
	authRoutes.PUT("/roles/:id", server.updateRole)
	authRoutes.DELETE("/roles/:id", server.softDeleteRole)

	authRoutes.POST("/role-users", server.createRoleUsers)
	authRoutes.GET("/role-users/:id", server.getRoleUserByUserID)
	authRoutes.GET("/role-users", server.listRoleUsers)
	authRoutes.PUT("/role-users/:id", server.updateRoleUsers)
	authRoutes.DELETE("/role-users/:id", server.softDeleteRoleUser)

	authRoutes.POST("/categories", server.createCategory)
	authRoutes.GET("/categories/:id", server.getCategory)
	authRoutes.GET("/categories", server.listCategory)
	authRoutes.PUT("/categories/:id", server.updateCategory)
	authRoutes.DELETE("/categories/:id", server.softDeleteCategory)

	authRoutes.POST("/partners", server.createPartner)
	authRoutes.GET("/partners/:id", server.getPartner)
	authRoutes.GET("/partners", server.listPartner)
	authRoutes.PUT("/partners/:id", server.updatePartner)
	authRoutes.DELETE("/partners/:id", server.softDeletePartner)

	authRoutes.POST("/providers", server.createProvider)
	authRoutes.GET("/providers/:id", server.getProvider)
	authRoutes.GET("/providers", server.listProvider)
	authRoutes.PUT("/providers/:id", server.updateProvider)
	authRoutes.DELETE("/providers/:id", server.softDeleteProvider)

	authRoutes.POST("/products", server.createProduct)
	authRoutes.GET("/products/:id", server.getProduct)
	authRoutes.GET("/products", server.listProduct)
	authRoutes.PUT("/products/:id", server.updateProduct)
	authRoutes.DELETE("/products/:id", server.softDeleteProduct)

	authRoutes.POST("/selling-prices", server.createSelling)
	authRoutes.GET("/selling-prices/:id", server.getSelling)
	authRoutes.GET("/selling-prices", server.listSelling)
	authRoutes.PUT("/selling-prices/:id", server.updateSelling)
	authRoutes.DELETE("/selling-prices/:id", server.softDeleteSelling)

	authRoutes.POST("/medias", server.createMedia)
	authRoutes.GET("/medias/one", server.getMedia)
	authRoutes.GET("/medias", server.listMedia)
	authRoutes.PUT("/medias/:id", server.updateMedia)
	authRoutes.DELETE("/medias/:id", server.softDeleteMedia)

	authRoutes.POST("/transactions", server.createTrx)
	authRoutes.GET("/transactions/:tx_id", server.getTrx)
	authRoutes.GET("/transactions", server.listTrx)
	authRoutes.PUT("/transactions/:tx_id", server.updateTrx)
	authRoutes.DELETE("/transactions/:id", server.softDeleteTrx)

	authRoutes.POST("/inquiry", server.inquiry)
	authRoutes.POST("/payment", server.payment)
	authRoutes.POST("/transactions/deposit", server.deposit)
	authRoutes.POST("/transactions/deposit-approve", server.depositApprove)

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
