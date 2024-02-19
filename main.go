package main

import (
	"comfyui_service/utils"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"net/http"
)

func main() {
	r := gin.Default()
	var prompt = utils.ReadWorkflowFile()
	var promptTemplate = utils.ReadPromptTemplate(utils.Normal)
	prompt.Process(func(key string, val utils.BaseNode) {
		val.UpdatePrompt(promptTemplate.PromptGroup)
		val.UpdateSampler(promptTemplate.Sampler)
		val.UpdateModel(promptTemplate.CheckPoint)
		val.UpdateOutputImage(promptTemplate.OutputImage)
	})
	prompt.AddLora(promptTemplate.Lora)
	utils.WriteWorkflowFile(prompt)
	if val, err := json.Marshal(prompt.Prompt); err == nil {
		println(string(val))
	}
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, prompt.Prompt)
		//c.JSON(http.StatusOK, gin.H{
		//	"data": prompt.Prompt,
		//})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
