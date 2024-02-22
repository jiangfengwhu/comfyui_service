package utils

import (
	"bytes"
	"fmt"
	"github.com/goccy/go-json"
	"io"
	"net/http"
	"os"
	"time"
)

type ServerConfig struct {
	ComfyHost        string `json:"comfy_host"`
	AddressUpdateUrl string `json:"address_update_url"`
	WorkflowDir      string `json:"workflow_dir"`
	TemplateDir      string `json:"template_dir"`
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
	res, _ := http.Get("https://api.ipify.org/")
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	http.Post(
		Config.AddressUpdateUrl,
		"application/json",
		bytes.NewBuffer([]byte(fmt.Sprintf(`{"addr":"%s","type":"t2i_addr"}`, string(body)))),
	)
}

func StartReportIpAddr() {
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
