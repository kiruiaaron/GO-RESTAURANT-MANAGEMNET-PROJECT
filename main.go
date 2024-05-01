package main

import (
	"os"

	"github.com/gin-gonic/gin"
	middleware "github.com/kiruiaaron/GO-RESTAURANT-MANAGEMNET-PROJECT/middlwares"
	"github.com/kiruiaaron/GO-RESTAURANT-MANAGEMNET-PROJECT/routes"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/kiruiaaron/GO-RESTAURANT-MANAGEMNET-PROJECT/database"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

func main(){
	port := os.Getenv("PORT")

	if port == ""{
		port = "8000"

	}

	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)
	router.Use(middleware.Authentication())

	routes.FoodRoutes(router)
	routes.MenuRoutes(router)
	routes.TableRoutes(router)
	routes.OrderRoutes(router)
	routes.OrderItemRoutes(router)
	routes.InvoiceRoutes(router)
	

	router.RU(":"+ port)



}