package utils

import (
	"fmt"
	"github.com/goccy/go-json"
	"math/rand"
	"os"
	"path/filepath"
)

type Meta struct {
	Title string `json:"title"`
}

type Prompt = map[string]BaseNode

type BaseNode struct {
	ClassType string                 `json:"class_type"`
	Inputs    map[string]interface{} `json:"inputs"`
	Meta      Meta                   `json:"_meta"`
}
type Workflow struct {
	Prompt      Prompt `json:"prompt"`
	StartLoraId string `json:"start_lora_id,omitempty"`
	ModelId     string `json:"model_id,omitempty"`
}

//======start node manipulation======

func (node *BaseNode) UpdateSeed() {
	if _, ok := node.Inputs["seed"]; ok {
		node.Inputs["seed"] = rand.Uint64()
	}
}
func (node *BaseNode) UpdateOutputImage(config OutputImage) {
	if _, ok := node.Inputs["batch_size"]; ok {
		node.Inputs["batch_size"] = config.BatchSize
		node.Inputs["width"] = config.Width
		node.Inputs["height"] = config.Height
	}
}
func (node *BaseNode) UpdateImagePrefix(prefix string) {
	if node.ClassType == "SaveImage" {
		node.Inputs["filename_prefix"] = prefix
	}
}
func (node *BaseNode) UpdateInputImage(image string) {
	if node.ClassType == "LoadImage" {
		node.Inputs["image"] = image
	}
}
func (node *BaseNode) UpdateSampler(config Sampler) {
	if _, ok := node.Inputs["cfg"]; ok {
		node.Inputs["cfg"] = config.Cfg
		node.Inputs["steps"] = config.Steps
		node.Inputs["sampler_name"] = config.SamplerName
		node.Inputs["scheduler"] = config.Scheduler
	}
}
func (node *BaseNode) UpdateModel(checkPoint string) {
	if node.ClassType == "CheckpointLoaderSimple" {
		node.Inputs["ckpt_name"] = checkPoint
	}
}
func (node *BaseNode) UpdatePrompt(promptGroup PromptGroup) {
	if "正向提示词" == node.Meta.Title {
		node.Inputs["text"] = promptGroup.Positive
	}
	if "负向提示词" == node.Meta.Title {
		node.Inputs["text"] = promptGroup.Negative
	}
}
func (node *BaseNode) GetLoraLocation() string {
	if "LoadStyleLoRA" == node.Meta.Title {
		return node.Inputs["model"].([]interface{})[0].(string)
	}
	return ""
}

//======end node manipulation======

//======start workflow manipulation======

func (workflow *Workflow) Process(handler func(string, BaseNode)) {
	for key, val := range workflow.Prompt {
		val.UpdateSeed()
		modelId := val.GetLoraLocation()
		if modelId != "" {
			workflow.StartLoraId = key
			workflow.ModelId = modelId
		}
		handler(key, val)
	}
}
func (workflow *Workflow) AddLora(loraConfig LoraConfig) {
	index := 0
	currentLoraId := workflow.StartLoraId
	lastLoraId := ""
	for loraName, strength := range loraConfig {
		index++
		newLora := workflow.NewLoraNode(loraName, strength)
		workflow.Prompt[currentLoraId] = newLora
		if lastLoraId != "" {
			workflow.Prompt[lastLoraId].Inputs["model"] = []interface{}{currentLoraId, 0}
			workflow.Prompt[lastLoraId].Inputs["clip"] = []interface{}{currentLoraId, 1}
		}
		lastLoraId = currentLoraId
		currentLoraId = fmt.Sprintf("lora_%d", index)
	}
}
func (workflow *Workflow) NewLoraNode(loraName string, strength float32) BaseNode {
	return BaseNode{
		Inputs: map[string]interface{}{
			"strength_model": strength,
			"strength_clip":  1,
			"lora_name":      loraName + ".safetensors",
			"model":          []interface{}{workflow.ModelId, 0},
			"clip":           []interface{}{workflow.ModelId, 1},
		},
		ClassType: "LoraLoader",
		Meta: Meta{
			Title: "LoadStyleLoRA",
		},
	}
}

// ======end workflow manipulation======

var WorkflowPool = map[string]Workflow{}

func ReadWorkflowFile(workflowType string) Workflow {
	if val, ok := WorkflowPool[workflowType]; ok {
		return val
	}
	file, err := os.ReadFile(filepath.Join(Config.WorkflowDir, workflowType+".json"))
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return Workflow{}
	}

	// 解析JSON数据到 map[string]interface{}
	var data map[string]BaseNode
	err = json.Unmarshal(file, &data)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return Workflow{}
	}
	WorkflowPool[workflowType] = Workflow{Prompt: data}
	return WorkflowPool[workflowType]
}

func WriteWorkflowFile(workflow Workflow) {
	if val, err := json.MarshalIndent(workflow.Prompt, "", "  "); err == nil {
		os.WriteFile("./workflows/tmp.json", val, 0644)
	}
}
