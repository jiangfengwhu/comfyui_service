package main

import (
	"comfyui_service/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/prompt", routes.GetPrompt)
	r.GET("/ws_secret_api", routes.LinkRemoteGPU)
	r.POST("/queue", routes.QueuePrompt)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
