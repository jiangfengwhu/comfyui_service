package routes

import (
	"comfyui_service/utils"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"net/http"
)

func GetPrompt(c *gin.Context) {
	templateId := c.Query("template_id")
	workflowId := c.Query("workflow_id")
	var prompt = utils.ReadWorkflowFile(workflowId)
	var promptTemplate = utils.ReadPromptTemplate(templateId)
	prompt.Process(func(key string, val utils.BaseNode) {
		val.UpdatePrompt(promptTemplate.PromptGroup)
		val.UpdateSampler(promptTemplate.Sampler)
		val.UpdateModel(promptTemplate.CheckPoint)
		val.UpdateOutputImage(promptTemplate.OutputImage)
	})
	prompt.AddLora(promptTemplate.Lora)
	//utils.WriteWorkflowFile(prompt)
	c.JSON(http.StatusOK, prompt.Prompt)
}

func queuePrompt(c *gin.Context) {
	templateId := c.Query("template_id")
	workflowId := c.Query("workflow_id")
	var prompt = utils.ReadWorkflowFile(workflowId)
	var promptTemplate = utils.ReadPromptTemplate(templateId)
	prompt.Process(func(key string, val utils.BaseNode) {
		val.UpdatePrompt(promptTemplate.PromptGroup)
		val.UpdateSampler(promptTemplate.Sampler)
		val.UpdateModel(promptTemplate.CheckPoint)
		val.UpdateOutputImage(promptTemplate.OutputImage)
	})
	prompt.AddLora(promptTemplate.Lora)
	promptJson, err := json.Marshal(prompt.Prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "JSON编码错误"})
		return
	}
	gpuConnection.WriteMessage(1, promptJson)
	//utils.WriteWorkflowFile(prompt)
	c.JSON(http.StatusOK, prompt.Prompt)
}
