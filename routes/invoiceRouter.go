package routes


import (
	"github.com/gin-gonic/gin"
	controller"github.com/kiruiaaron/GO-RESTAURANT-MANAGEMNET-PROJECT/controllers"

)
func InvoiceRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.GET("/invoices", controller.GetInvoices())
	incomingRoutes.GET("/invoices/:invoice_id",controller.GetInvoice())
	incomingRoutes.POST("/invoices", controller.CreateInvoices())
	incomingRoutes.PATCH("/invoices/invoice_id", controller.UpdateInvoice())

}