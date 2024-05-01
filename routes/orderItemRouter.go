package routes


import (
	"github.com/gin-gonic/gin"
	controller"github.com/kiruiaaron/GO-RESTAURANT-MANAGEMNET-PROJECT/controllers"

)
func OrderItemRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.GET("/orderItems", controller.GetOrderItems())
	incomingRoutes.GET("/orderItems/:orderItem_id",controller.GetOrderItem())
	incomingRoutes.GET("/orderItems-order/:order_id",controller.GetOrderItemsByOrder())
	incomingRoutes.POST("/orderItems", controller.CreateOrderItems())
	incomingRoutes.PATCH("/orderItems/orderItem_id", controller.UpdateOrderItem())

}