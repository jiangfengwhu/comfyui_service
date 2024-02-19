package utils

import (
	"fmt"
	"github.com/goccy/go-json"
	"os"
)

type TemplateType string
type LoraConfig map[string]float32

const (
	SDFacePortrait TemplateType = "SDFacePortrait"
	Normal         TemplateType = "基础版"
)

type Sampler struct {
	Steps       int     `json:"steps"`
	Cfg         float32 `json:"cfg"`
	SamplerName string  `json:"sampler_name"`
	Scheduler   string  `json:"scheduler"`
}
type PromptGroup struct {
	Positive string `json:"positive"`
	Negative string `json:"negative"`
}
type OutputImage struct {
	Width     int `json:"width"`
	Height    int `json:"height"`
	BatchSize int `json:"batch_size"`
}
type PromptTemplate struct {
	OutputImage OutputImage  `json:"output_image"`
	Sampler     Sampler      `json:"sampler"`
	Type        TemplateType `json:"type"`
	CheckPoint  string       `json:"check_point"`
	Lora        LoraConfig   `json:"lora"`
	PromptGroup PromptGroup  `json:"prompt_group"`
	Name        string       `json:"name,omitempty"`
}

var TemplatePool = map[TemplateType]PromptTemplate{}

func ReadPromptTemplate(tmpType TemplateType) PromptTemplate {
	if val, ok := TemplatePool[tmpType]; ok {
		return val
	}
	file, err := os.ReadFile(fmt.Sprintf("./templates/%s.json", tmpType))
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return PromptTemplate{}
	}

	// 解析JSON数据到 map[string]interface{}
	var data PromptTemplate
	err = json.Unmarshal(file, &data)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return PromptTemplate{}
	}
	TemplatePool[tmpType] = data
	return data
}
