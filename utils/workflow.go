package utils

import (
	"fmt"
	"os"
)

type Meta struct {
	Title string `json:"title"`
}
type BaseNode struct {
	ClassType string                 `json:"class_type"`
	Inputs    map[string]interface{} `json:"inputs"`
	Mata      Meta                   `json:"_meta"`
}
type KSamplerNode struct {
	BaseNode
}

func (node *KSamplerNode) UpdateSeed() {
	node.Inputs["seed"] = 565011220517970
	node.Inputs["model"] = []interface{}{"1", 0}
}

func CreateCheckPointLoader(modelPath string) BaseNode {
	return BaseNode{
		ClassType: "CheckpointLoaderSimple",
		Inputs: map[string]interface{}{
			"ckpt_name": modelPath,
		},
		Mata: Meta{
			Title: "Load Checkpoint",
		},
	}
}

func CreateKSampler() KSamplerNode {
	return KSamplerNode{
		BaseNode: BaseNode{
			ClassType: "KSampler",
			Inputs: map[string]interface{}{
				"seed":         565011220517970,
				"steps":        30,
				"cfg":          7,
				"sampler_name": "dpmpp_2m_sde",
				"scheduler":    "karras",
				"denoise":      1,
				"model":        []interface{}{},
				"positive":     []interface{}{},
				"negative":     []interface{}{},
				"latent_image": []interface{}{},
			},
			Mata: Meta{
				Title: "Sample",
			},
		},
	}
}

func ReadWorkflowFile() string {
	file, err := os.ReadFile("./workflows/ip.json")
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return ""
	}

	// 解析JSON数据到 map[string]interface{}
	return string(file)
}
