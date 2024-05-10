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
)


type InvoiceViewFormat struct{
    Invoice_id         string
	Payment_method     string
	Order_id            string
	Payment_status      *string
	Payment_due        interface{}
	Table_number       interface{}
	Payment_due_date    time.Time
	Order_details      interface{}
}

var invoiceCollection *mongo.Collection = database.OpenCollection(database.Client, "invoice")


func GetInvoices() gin.HandlerFunc{
	return func (c *gin.Context)  {
		var ctx,cancel = context.WithTimeout(context.Background(), 100*time.Second)

		result, err := invoiceCollection.Find(context.TODO(),bson.M{})
		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"error occurred while creating invoice"})
			return
		}
		var allInvoices []bson.M
		if err - result.All(ctx, &allInvoices); err !=nil{
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allInvoices)
	}
}


func GetInvoice() gin.HandlerFunc{
	return func (c *gin.Context)  {
		var ctx,cancel = context.WithTimeout(context.Background(), 100*time.Second)
        invoiceId := c.Param("invoice_id")

		var invoice models.Invoice
		err := invoiceCollection.FindOne(ctx, bson.M{"invoice_id":invoiceId}).Decode(&invoice)
		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"error occurred while listing invoice item"})
			return
		}

		var invoiceView InvoiceViewFormat

		allOrderItems, err := ItemsByOrder(invoice.Order_id)
		invoiceView.Order_id = invoice.Order_id
		invoiceView.Payment_due_date = invoice.Payment_Due_Date

		invoiceView.Payment_method = "null"
		if invoice.Payment_Method != nil{
			invoiceView.Payment_method = *invoice.Payment_Method 
		}

		invoiceView.Invoice_id = invoice.Invoice_id
		invoiceView.Payment_status = *&invoice.Payment_Status
		invoiceView.Payment_due = allOrderItems[0]["payment_due"]
		invoiceView.Table_number = allOrderItems[0]["table_number"]
		invoiceView.Order_details = allOrderItems[0]["order_items"]

		c.JSON(http.StatusOK, invoiceView)
	}
}


func CreateInvoice() gin.HandlerFunc{
	return func (c *gin.Context)  {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var invoice models.Invoice

		if err := c.BindJSON(&invoice); err!= nil{
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
			return
		}

		var order models.Order

		err := orderCollection.FindOne(ctx, bson.M{"order_id":order.Order_id}).Decode(&order)
			defer cancel()
			if err != nil{
				msg := fmt.Sprintf("message: Menu was not found")
				c.JSON(http.StatusInternalServerError, gin.H{"error":msg})
				return
		} 
		status := "PENDING"
		if invoice.Payment_Status == nil{
			invoice.Payment_Status = &status
		}  
		invoice.Payment_Due_Date, _ =  time.Parse(time.RFC3339, time.Now().AddDate(0,0,1).Format(time.RFC3339))
		invoice.Created_at, _ =time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		invoice.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		invoice.ID = primitive.NewObjectID()
		invoice.Invoice_id = invoice.ID.Hex()

		validationErr := validate.Struct(invoice)
		if validationErr != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error":validationErr.Error()})
			return
		}

		result, insertErr  := invoiceCollection.InsertOne(ctx, invoice)
		if insertErr != nil{
            msg := fmt.Sprintf("invoice item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error":msg})
			return
		}

		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}


func UpdateInvoice() gin.HandlerFunc{
	return func (c *gin.Context)  {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var invoice models.Invoice
		invoiceId := c.Param("invoice_id")
		if err := c.BindJSON(&invoice); err!= nil{
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
			return
		}

		filter := bson.M{"invoice_id":invoiceId}
		var updatedObj primitive.D

		if invoice.Payment_Method != nil{
            updatedObj = append(updatedObj, bson.E{"payment_method", invoice.Payment_Method})
		}

		if invoice.Payment_Status != nil{
            updatedObj = append(updatedObj, bson.E{"payment_status", invoice.Payment_Status})
		}

		invoice.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updatedObj = append(updatedObj, bson.E{"updated_at", invoice.Updated_at})

		upsert := true

		opt := options.UpdateOne{
			Upsert : &upsert,
		}
        
		
		result, err := invoiceCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$set", updatedObj},
			},
			&opt,
		)


		if err != nil{
			msg:= fmt.Sprintf("invoice item update failed")
			c.JSON(http.StatusInternalServerError, gin.H{"error":msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}