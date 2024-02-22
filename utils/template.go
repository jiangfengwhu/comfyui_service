package utils

import (
	"fmt"
	"github.com/goccy/go-json"
	"os"
	"path/filepath"
)

type LoraConfig map[string]float32

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
	OutputImage OutputImage `json:"output_image"`
	Sampler     Sampler     `json:"sampler"`
	Type        string      `json:"type"`
	CheckPoint  string      `json:"check_point"`
	Lora        LoraConfig  `json:"lora"`
	PromptGroup PromptGroup `json:"prompt_group"`
	Name        string      `json:"name,omitempty"`
}

var TemplatePool = map[string]PromptTemplate{}

func ReadPromptTemplate(tmpType string) PromptTemplate {
	if val, ok := TemplatePool[tmpType]; ok {
		return val
	}
	file, err := os.ReadFile(filepath.Join(Config.TemplateDir, tmpType+".json"))
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
