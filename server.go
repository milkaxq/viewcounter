package main

import (
	"database/sql"
	"net/http"
	"os"
	"time"

	"viewscounter/config"
	"viewscounter/controller"
	"viewscounter/helper"
	"viewscounter/repository"
	"viewscounter/service"
	"viewscounter/ws"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
)

var (
	db                *sql.DB                      = config.SetupDatabaseConnection()
	productRepository repository.ProductRepository = repository.NewProductRepository(db)
	productService    service.ProductService       = service.NewProductService(productRepository)
	productController controller.ProductController = controller.NewProductController(productService)
)

func main() {
	h := ws.NewHub()
	go h.Run()
	//Call EnvParser
	config.EnvParser()

	mode := os.Getenv("GIN_MODE")
	gin.SetMode(mode)
	//Instance of Gin Framework
	r := gin.Default()
	//Cors
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, //Change
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	viewCounterRoute := os.Getenv("VIEW_COUNTER_ROUTE")
	prefixRoute := os.Getenv("PREFIX_ROUTE")
	productFindRoute := os.Getenv("PRODUCT_SEARCH_ROUTE")
	ws_orderInvRoute := os.Getenv("WS_ORDER_INV_ROUTE")
	ws_smsRegisterRoute := os.Getenv("WS_SMS_REGISTER_ROUTE")
	api_orderInvRoute := os.Getenv("API_ORDER_INV_ROUTE")
	api_smsRegisterRoute := os.Getenv("API_SMS_REGISTER_ROUTE")

	//Group routes of same origin
	productRoutes := r.Group(prefixRoute)
	{
		productRoutes.GET(viewCounterRoute, productController.FilterByIpAndUserAgent)
		productRoutes.GET(productFindRoute, productController.FindProduct)
	}
	r.POST(api_orderInvRoute, func(c *gin.Context) {
		token := c.GetHeader("x-access-token")
		if token != os.Getenv("SHA_KEY") {
			res := helper.BuildErrorResponse("Failed to get token", "error", helper.EmptyObj{})
			c.JSON(http.StatusUnauthorized, res)
		} else {
			jsonData, err := c.GetRawData()
			if err != nil {
				res := helper.BuildErrorResponse("Failed Get Json", err.Error(), helper.EmptyObj{})
				c.JSON(http.StatusBadRequest, res)
			}
			m := ws.Message{Data: jsonData, Room: "orderRevision"}
			h.Broadcast <- m

			res := helper.BuildResponse(true, "Sended", helper.EmptyObj{})
			c.JSON(http.StatusOK, res)
		}
	})

	r.POST(api_smsRegisterRoute, func(c *gin.Context) {
		token := c.GetHeader("x-access-token")
		if token != os.Getenv("SHA_KEY") {
			res := helper.BuildErrorResponse("Failed to get token", "error", helper.EmptyObj{})
			c.JSON(http.StatusUnauthorized, res)
		} else {
			jsonData, err := c.GetRawData()
			if err != nil {
				res := helper.BuildErrorResponse("Failed Get Json", err.Error(), helper.EmptyObj{})
				c.JSON(http.StatusBadRequest, res)
			}
			m := ws.Message{Data: jsonData, Room: "smsRegister"}
			h.Broadcast <- m

			res := helper.BuildResponse(true, "Sended", helper.EmptyObj{})
			c.JSON(http.StatusOK, res)
		}
	})
	//WS STARTS HERE
	r.GET("/ws/*name", func(c *gin.Context) {
		name := c.Param("name")
		if name == ws_orderInvRoute {
			ws.ServeWs(h, c.Writer, c.Request, "orderRevision")
		} else if name == ws_smsRegisterRoute {
			ws.ServeSMS(h, c.Writer, c.Request, "smsRegister")
		}
	})

	//Cron - Periodicly call function
	s := gocron.NewScheduler(time.UTC)
	s.Every(os.Getenv("VIEWS_UPDATE_TIME")).Do(productService.InsertViews)
	s.Every(os.Getenv("DB_PRODUCT_SYNC_TIME")).Do(productService.BleveSync)
	s.StartAsync()
	r.Run()
}
