package main

import (
	"comfyui_service/routes"
	"comfyui_service/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	utils.InitConfig()
	utils.StartReportIpAddr()
	r := gin.Default()
	r.GET("/prompt", routes.GetPrompt)
	r.POST("/queue_prompt", routes.QueuePrompt)
	r.GET("/alive", routes.ALive)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
