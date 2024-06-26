package routes


import (
	"github.com/gin-gonic/gin"
	controller"github.com/kiruiaaron/GO-RESTAURANT-MANAGEMNET-PROJECT/controllers"

)
func TableRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.GET("/tables", controller.GetTables())
	incomingRoutes.GET("/tables/:table_id",controller.GetTable())
	incomingRoutes.POST("/tables", controller.CreateTables())
	incomingRoutes.PATCH("/tables/table_id", controller.UpdateTables())

}