package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"viewscounter/config"
	"viewscounter/controller"
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
	jwtService        service.JWTService           = service.NewJWTService()
	productController controller.ProductController = controller.NewProductController(productService)
)

func main() {
	h := ws.NewHub()
	go h.Run()
	//Call EnvParser
	config.EnvParser()
	//Instance of Gin Framework
	r := gin.Default()

	mode := os.Getenv("GIN_MODE")
	gin.SetMode(mode)
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
	productRoute := os.Getenv("PRODUCT_ROUTE")
	findRoute := os.Getenv("SEARCH_ROUTE")
	//Group routes of same origin

	productRoutes := r.Group(productRoute)
	{
		productRoutes.GET(viewCounterRoute, productController.FilterByIpAndUserAgent)
		productRoutes.GET(findRoute, productController.FindProduct)
	}
	r.POST("/order-inv/", func(c *gin.Context) {
		jsonData, err := c.GetRawData()
		if err != nil {
			log.Println(err)
		}
		m := ws.Message{Data: jsonData, Room: "orderRevision"}
		h.Broadcast <- m
	})
	//WS STARTS HERE
	r.GET("/ws/", func(c *gin.Context) {
		ws.ServeWs(h, c.Writer, c.Request, "orderRevision")
	})

	//Cron - Periodicly call function
	s := gocron.NewScheduler(time.UTC)
	s.Every(os.Getenv("VIEWS_UPDATE_TIME")).Do(productService.InsertViews)
	s.Every(os.Getenv("DB_PRODUCT_SYNC_TIME")).Do(productService.BleveSync)
	s.StartAsync()
	r.Run()
}
