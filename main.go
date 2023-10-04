package main

import (
	"log"
	"os"

	"github.com/MrBooi/ecommerce-cart/controllers"
	"github.com/MrBooi/ecommerce-cart/database"

	middleware "github.com/MrBooi/ecommerce-cart/middleware"
	"github.com/MrBooi/ecommerce-cart/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	app := controllers.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, "Users"))

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(middleware.Authentication())

	routes.UserRoutes(router)

	router.POST("/addtocart", app.AddToCart())
	router.DELETE("/removeitem", app.RemoveItem())
	router.POST("/cartcheckout", app.BuyFromCart())
	router.POST("/instantbuy", app.InstantBuy())

	log.Fatal(router.Run(":" + port))
}
