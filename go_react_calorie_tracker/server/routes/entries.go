package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kaung-minkhant/go_projs/go_react_calorie_tracker/controllers"
)

func setupEntries(engine *gin.Engine) {
  engine.POST("/entry/create", controllers.AddEntry)
  engine.GET("/entry", controllers.GetEntries)
  engine.GET("/entry/:id", controllers.GetEntryById)
  engine.GET("/ingredient/:ingredient", controllers.GetEntriesByIngredient)
  engine.PUT("/entry/:id", controllers.UpdateEntry)
  engine.PUT("/ingredient/:ingredient", controllers.UpdateIngredient)
  engine.DELETE("/entry/:id", controllers.DeleteEntry)
}
