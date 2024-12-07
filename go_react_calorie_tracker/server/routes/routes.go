package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
  engine := gin.New()
  initMiddleWares(engine)
  setupEntries(engine)
  return engine
}

func initMiddleWares(engine *gin.Engine) {
  engine.Use(gin.Logger())
  engine.Use(cors.Default())
}
