package routes

import (
	"comfyui_service/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type QueuePromptReq struct {
	Type       string            `json:"type" binding:"required"`
	TemplateId string            `json:"template_id" binding:"required"`
	WorkflowId string            `json:"workflow_id" binding:"required"`
	Images     map[string]string `json:"images,omitempty"`
}

type ImageUploadReq struct {
	ImgBase64 string `json:"img_base64"`
}

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
