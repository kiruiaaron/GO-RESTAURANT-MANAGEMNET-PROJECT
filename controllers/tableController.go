package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kiruiaaron/GO-RESTAURANT-MANAGEMNET-PROJECT/database"
	"github.com/kiruiaaron/GO-RESTAURANT-MANAGEMNET-PROJECT/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)


var tableCollection *mongo.Collection = database.OpenCollection(database.Client, "table")

func GetTables() *gin.HandlerFunc{
	return func (c *gin.Context)  {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		result, err := tableCollection.Find(context.TODO(), bson.M{})

		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"error occurred while listing table items"})
			return
		}

		var allTables []bson.M
		if err = result.All(ctx, &allTables); err != nil{
			log.Fatal(err)
			return
		}
		c.JSON(http.StatusOK, allTables)
	}
}

func GetTable() *gin.HandlerFunc{
	return func (c *gin.Context)  {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		tableId := c.Param("table_id")
		var tableItem models.Table

		err := tableCollection.FindOne(ctx, bson.M{"table_id":tableId}).Decode(&tableItem)
		defer cancel()

		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"error occurred while listing order item"})
			return
		}
		c.JSON(http.StatusOK, tableItem)
	}
}
func CreateTable() *gin.HandlerFunc{
	return func (c *gin.Context)  {
		
	}
}
func UpdateTable() *gin.HandlerFunc{
	return func (c *gin.Context)  {
		
	}
}