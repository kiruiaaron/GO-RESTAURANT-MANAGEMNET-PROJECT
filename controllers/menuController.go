package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kiruiaaron/GO-RESTAURANT-MANAGEMNET-PROJECT/database"
	"github.com/kiruiaaron/GO-RESTAURANT-MANAGEMNET-PROJECT/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/mongo/driver/mongocrypt/options"
	"golang.org/x/tools/go/analysis/passes/defers"
)

var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")

func GetMenus() *gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		result, err := menuCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing the menu"})
		}
		var allMenu []bson.M
		if err = result.All(ctx, &allMenu); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allMenu)
	}
}

func GetMenu() *gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		menuId := c.Param("menu_id")
		var menu models.Menu

		err := foodCollection.FindOne(ctx, bson.M{"meu_id": menuId}).Decode(&menu)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while fetching menu"})
		}
		c.JSON(http.StatusOK, menu)

	}
}
func CreateMenu() *gin.HandlerFunc {
	return func(c *gin.Context) {
		var menu models.Menu
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(menu)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		menu.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.ID = primitive.NewObjectID()
		menu.Menu_id = menu.ID.Hex()

		results, insertErr := menuCollection.InsertOne(ctx, menu)
		if insertErr != nil {
			msg := fmt.Sprintf("Menu item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, results)
	}
}

func inTimeSpan(start, end, check time.Time)bool{
	return start.After(time.Now()) && end.After(start)
}



func UpdateMenu() *gin.HandlerFunc {
	return func(c *gin.Context) {
        var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var menu models.Menu

		if err := c.BindJSON(&menu); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
			return
		}

		menuId := c.Param("menu_id")
		filter := bson.M{"menu_id":menuId}

		var updatedObj primitive.D

		if menu.Start_Date != nil && menu.End_Date != nil{
			if !inTimeSpan(*menu.Start_Date, *menu.End_Date, time.Now()){
				msg := "Kindly retype the time"
				c.JSON(http.StatusInternalServerError, gin.H{"error":msg})
				defer cancel()
				return
			}

			updatedObj = append(updatedObj, bson.E{"start_date",menu.Start_Date})
			updatedObj = append(updatedObj, bson.E{"end_date", menu.End_Date})

			if menu.Name != ""{
				updatedObj = append(updatedObj, bson.E{"name", menu.Name})
			}

			if menu.Category != ""{
				updatedObj = append(updatedObj, bson.E{"category", menu.Category})
			}

			menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			updatedObj = append(updatedObj, bson.E{"updated_at", menu.Updated_at})

			upsert := true

			opt := options.UpdateOptions{
				Upsert:&upsert,
                     
			}

			result, err := menuCollection.UpdateOne(
				ctx,
				filter,
				bson.D{
					{"$set", updatedObj},
				},
				&opt,
			)
			if err != nil{
				msg:= "Menu update failed"
				c.JSON(http.StatusInternalServerError, gin.H{"error":msg})

			}
			defer cancel()
			c.JSON(http.StatusOK,result)
		}

	}

}