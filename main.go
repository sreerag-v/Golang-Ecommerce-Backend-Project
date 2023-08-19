package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sreerag_v/Ecom/database"
	"github.com/sreerag_v/Ecom/routes"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	database.InitDB()

	// gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")
	routes.UserRoutes(router)
	routes.AdminRoutes(router)
	router.Run()
}
