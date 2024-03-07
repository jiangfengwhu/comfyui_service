package utils

import (
	"fmt"
	"github.com/goccy/go-json"
	"os"
	"path/filepath"
	"strings"
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
	Desc        string      `json:"desc,omitempty"`
}
type ImageItem struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}

var TemplatePool = map[string]PromptTemplate{}

func addToTemplatePool(tmpType string) error {
	file, err := os.ReadFile(filepath.Join(Config.TemplateDir, tmpType+".json"))
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return err
	}

	// 解析JSON数据到 map[string]interface{}
	var data PromptTemplate
	err = json.Unmarshal(file, &data)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return err
	}
	TemplatePool[tmpType] = data
	return nil
}
func ReadPromptTemplate(tmpType string) PromptTemplate {
	if val, ok := TemplatePool[tmpType]; ok {
		return val
	}
	err := addToTemplatePool(tmpType)
	if err != nil {
		return PromptTemplate{}
	}
	return TemplatePool[tmpType]
}

func UpdateTemplatePool() {
	files, err := os.ReadDir(Config.TemplateDir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileName := file.Name()
		if filepath.Ext(fileName) != ".json" {
			continue
		}
		err := addToTemplatePool(strings.TrimSuffix(fileName, ".json"))
		if err != nil {
			fmt.Println("Error reading file:", err)
		}
	}
}

func GetAllTemplateId() []interface{} {
	keys := make([]interface{}, 0, len(TemplatePool))
	for key, val := range TemplatePool {
		keys = append(keys, map[string]interface{}{"id": key, "desc": val.Desc})
	}
	return keys
}

var homeList []ImageItem

func GetHomeList(refresh bool) []ImageItem {
	if homeList == nil || refresh {
		homeList = []ImageItem{}
		files, err := os.ReadDir(Config.HomeImgDir)
		if err != nil {
			fmt.Println("Error reading directory:", err)
			return []ImageItem{}
		}
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			fileName := file.Name()
			chunks := strings.Split(fileName, "_")
			homeList = append(homeList, ImageItem{Id: chunks[0], Url: fileName})
		}
	}
	return homeList
}
