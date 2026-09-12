package api

import (
	"Mav_backend/api/handlers"
	"Mav_backend/api/middleware"
	"Mav_backend/database"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	pool, err := database.DatabaseInit()
	if err != nil {
		panic(err)
	}

	h := handlers.NewDbPool(pool)
	p := r.Group("/")

	p.GET("/health_data", h.GetRecent)
	p.GET("/Word_of_the_day", h.GetWordOfTheDay)
	rest := r.Group("/r")
	rest.Use(middleware.CheckTokenMiddleware())

	rest.POST("/health_data", h.CreateHealthData)
	rest.POST("/Word_of_the_day", h.CreateWordOfTheDay)

	return r
}
