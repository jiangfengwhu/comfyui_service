package routes

import (
	"bytes"
	"comfyui_service/db"
	"comfyui_service/model"
	"comfyui_service/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"go.mongodb.org/mongo-driver/bson"
	"io"
	"log"
	"net/http"
)

type QueuePromptReq struct {
	Type       string            `json:"type" binding:"required"`
	TemplateId string            `json:"template_id" binding:"required"`
	Images     map[string]string `json:"images,omitempty"`
	HomeMode   bool              `json:"home_mode,omitempty"`
}

type ImageUploadReq struct {
	ImgBase64 string `json:"img_base64"`
}

func queuePrompt(data utils.Prompt) (string, error) {
	jsonData, err := json.Marshal(map[string]interface{}{"prompt": data})
	if err != nil {
		log.Println("marshal err:", err)
		return "", err
	}
	url := fmt.Sprintf("%s/prompt", utils.Config.ComfyHost)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("post err:", err)
		return "", err
	}
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("read err:", err)
		return "", err
	}
	var result map[string]interface{}
	err = json.Unmarshal(respBytes, &result)
	if err != nil {
		log.Println("unmarshal err:", err)
		return "", err
	}
	return result["prompt_id"].(string), nil
}

func GetPrompt(c *gin.Context) {
	templateId := c.Query("template_id")
	var promptTemplate = utils.ReadPromptTemplate(templateId)
	workflowId := promptTemplate.Type
	var prompt = utils.ReadWorkflowFile(workflowId)
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
	var promptTemplate = utils.ReadPromptTemplate(req.TemplateId)
	workflowId := promptTemplate.Type
	var prompt = utils.ReadWorkflowFile(workflowId)
	outputPrefix := utils.GetUUID()
	prompt.Process(func(key string, val utils.BaseNode) {
		val.UpdatePrompt(promptTemplate.PromptGroup)
		val.UpdateSampler(promptTemplate.Sampler)
		val.UpdateModel(promptTemplate.CheckPoint)
		val.UpdateOutputImage(promptTemplate.OutputImage)
		if req.HomeMode {
			prefix := fmt.Sprintf("%s_%d_%d", req.TemplateId, promptTemplate.OutputImage.Width, promptTemplate.OutputImage.Height)
			val.UpdateImagePrefix(prefix, "home", "webp")
		} else {
			val.UpdateImagePrefix(outputPrefix, "", "jpg")
		}
		if img, ok := req.Images[key]; ok {
			val.UpdateInputImage(img)
		}
	})
	prompt.AddLora(promptTemplate.Lora)
	promptId, err := queuePrompt(prompt.Prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "加入队列失败" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.Response{Code: 0, Msg: "success", Data: map[string]string{"prompt_id": promptId, "output_prefix": outputPrefix}})
}

func GetTemplates(c *gin.Context) {
	c.JSON(http.StatusOK, model.Response{Code: 0, Data: utils.GetAllTemplateId()})
}

func GetHomeList(c *gin.Context) {
	refresh := c.Query("refresh")
	c.JSON(http.StatusOK, model.Response{Code: 0, Data: utils.GetHomeList(refresh == "true")})
}

func UpdateTemplate(c *gin.Context) {
	utils.UpdateTemplatePool()
	c.JSON(http.StatusOK, model.Response{Code: 0, Msg: "success"})
}

func UpdateWorkflow(c *gin.Context) {
	utils.UpdateWorkflowPool()
	c.JSON(http.StatusOK, model.Response{Code: 0, Msg: "success"})
}

func MiniQueuePrompt(c *gin.Context) {
	var req QueuePromptReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := c.MustGet("user").(db.User)
	if user.Tickets < 1 {
		c.JSON(http.StatusOK, model.Response{
			Code: -1,
			Msg:  "额度不足",
		})
		return
	}
	var promptTemplate = utils.ReadPromptTemplate(req.TemplateId)
	workflowId := promptTemplate.Type
	var prompt = utils.ReadWorkflowFile(workflowId)
	outputPrefix := utils.GetUUID()
	prompt.Process(func(key string, val utils.BaseNode) {
		val.UpdatePrompt(promptTemplate.PromptGroup)
		val.UpdateSampler(promptTemplate.Sampler)
		val.UpdateModel(promptTemplate.CheckPoint)
		val.UpdateOutputImage(promptTemplate.OutputImage)
		val.UpdateImagePrefix(outputPrefix, "", "jpg")
		if img, ok := req.Images[key]; ok {
			val.UpdateInputImage(img)
		}
	})
	prompt.AddLora(promptTemplate.Lora)
	promptId, err := queuePrompt(prompt.Prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Code: -1, Msg: "加入队列失败" + err.Error()})
		return
	}
	filter := bson.M{"_id": user.Id}
	updater := bson.M{"$set": bson.M{"tickets": user.Tickets - 1}, "$push": bson.M{"history": outputPrefix}}
	db.UpdateUserOne(filter, updater)
	c.JSON(http.StatusOK, model.Response{Code: 0, Msg: "success", Data: map[string]string{"prompt_id": promptId, "output_prefix": outputPrefix}})
}
