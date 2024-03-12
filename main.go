package main

import (
	"comfyui_service/db"
	"comfyui_service/routes"
	"comfyui_service/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	utils.InitConfig()
	utils.StartReportIpAddr()
	utils.UpdateTemplatePool()
	db.Init()
	defer db.CloseDB()

	r := gin.Default()
	r.GET("/prompt", routes.GetPrompt)
	r.POST("/queue_prompt", routes.QueuePrompt)
	r.GET("/alive", routes.ALive)
	r.GET("/templates", routes.GetTemplates)
	r.GET("/refresh_template", routes.UpdateTemplate)
	r.GET("/refresh_workflow", routes.UpdateWorkflow)
	r.GET("/home", routes.GetHomeList)
	r.GET("/login", routes.Login)
	miniAuth := r.Group("/mini")
	miniAuth.Use(db.AuthRequired)
	{
		miniAuth.GET("/alive", routes.ALive)
	}
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
