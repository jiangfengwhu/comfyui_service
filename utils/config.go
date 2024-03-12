package utils

import (
	"bytes"
	"fmt"
	"github.com/goccy/go-json"
	"net/http"
	"os"
	"time"
)

type ServerConfig struct {
	ComfyHost        string `json:"comfy_host"`
	AddressUpdateUrl string `json:"address_update_url"`
	WorkflowDir      string `json:"workflow_dir"`
	TemplateDir      string `json:"template_dir"`
	HomeImgDir       string `json:"home_img_dir"`
	MongoUrl         string `json:"mongo_url"`
	APPID            string `json:"appid"`
	APPSecret        string `json:"app_secret"`
}

func (config *ServerConfig) String() string {
	return config.ComfyHost
}

var Config ServerConfig

func InitConfig() {
	fileData, err := os.ReadFile("./config.json")
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}
	err = json.Unmarshal(fileData, &Config)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}
	fmt.Println("Config loaded successfully: ", Config)
}

func ReportIpAddr() {
	resp, err := http.Post(
		Config.AddressUpdateUrl,
		"application/json",
		bytes.NewBuffer([]byte(`{"type":"t2i_addr"}`)),
	)
	if err != nil {
		fmt.Println("Error Report Ip Addr:", err)
		return
	}
	defer resp.Body.Close()
}

func StartReportIpAddr() {
	ReportIpAddr()
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				ReportIpAddr()
			}
		}
	}()
}
