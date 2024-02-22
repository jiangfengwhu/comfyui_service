package routes

import (
	"bytes"
	"comfyui_service/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"log"
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

func queuePrompt(data utils.Prompt) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println("marshal err:", err)
		return err
	}
	url := fmt.Sprintf("%s/prompt", utils.Config.ComfyHost)
	_, err = http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	return err
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
	err := queuePrompt(prompt.Prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "加入队列失败" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "已加入队列"})
}
