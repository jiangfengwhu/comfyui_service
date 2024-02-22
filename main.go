package main

import (
	"comfyui_service/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/prompt", routes.GetPrompt)
	r.GET("/alive", routes.ALive)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
