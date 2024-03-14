package main

import (
	"comfyui_service/db"
	"comfyui_service/routes"
	"comfyui_service/utils"
	"comfyui_service/ws_client"
	"github.com/gin-gonic/gin"
)

func main() {
	utils.InitConfig()
	utils.StartReportIpAddr()
	utils.UpdateTemplatePool()
	db.Init()
	defer db.CloseDB()
	ws_client.InitWs()

	r := gin.Default()
	r.GET("/prompt", routes.GetPrompt)
	r.POST("/queue_prompt", routes.QueuePrompt)
	r.GET("/alive", routes.ALive)
	r.GET("/templates", routes.GetTemplates)
	r.GET("/refresh_template", routes.UpdateTemplate)
	r.GET("/refresh_workflow", routes.UpdateWorkflow)
	r.GET("/home", routes.GetHomeList)
	r.GET("/login_wx", routes.Login)
	miniAuth := r.Group("/wx")
	miniAuth.Use(db.AuthRequired)
	{
		miniAuth.POST("/queue_prompt", routes.MiniQueuePrompt)
		miniAuth.GET("/user_info", routes.UserInfo)
		miniAuth.POST("/upload", routes.UploadFile)
		miniAuth.GET("/my_gallery", routes.MyGallery)
	}
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
