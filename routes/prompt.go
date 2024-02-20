package routes

import (
	"comfyui_service/utils"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
	"net/http"
)

type QueuePromptReq struct {
	Type       string            `json:"type" binding:"required"`
	TemplateId string            `json:"template_id" binding:"required"`
	WorkflowId string            `json:"workflow_id" binding:"required"`
	Images     map[string]string `json:"images,omitempty"`
}

type ToGPU struct {
	Type   string                    `json:"type"`
	Prompt map[string]utils.BaseNode `json:"prompt"`
	Images map[string]string         `json:"images,omitempty"`
}
type ImageUploadReq struct {
	ImgBase64 string `json:"img_base64"`
}
type FromGPU struct {
	Type    string `json:"type"`
	Content string `json:"content,omitempty"`
	Image   []byte `json:"image,omitempty"`
	Id      int16  `json:"id,omitempty"`
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

func QueuePrompt(c *gin.Context) {
	var req QueuePromptReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var prompt = utils.ReadWorkflowFile(req.WorkflowId)
	var promptTemplate = utils.ReadPromptTemplate(req.TemplateId)
	prompt.Process(func(key string, val utils.BaseNode) {
		val.UpdatePrompt(promptTemplate.PromptGroup)
		val.UpdateSampler(promptTemplate.Sampler)
		val.UpdateModel(promptTemplate.CheckPoint)
		val.UpdateOutputImage(promptTemplate.OutputImage)
	})
	prompt.AddLora(promptTemplate.Lora)
	gpuMessage := ToGPU{
		Type:   "prompt",
		Prompt: prompt.Prompt,
		Images: req.Images,
	}
	gpuMsgText, err := json.Marshal(gpuMessage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "JSON编码错误"})
		return
	}
	err = gpuConnection.WriteMessage(websocket.TextMessage, gpuMsgText)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "GPU连接错误"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "已加入队列"})
}
