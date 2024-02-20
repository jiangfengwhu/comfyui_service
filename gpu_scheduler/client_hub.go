package main

import (
	"bytes"
	"comfyui_service/routes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/goccy/go-json"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

type Config struct {
	Addr           string `json:"addr"`
	InputImageDir  string `json:"input_image_dir"`
	OutputImageDir string `json:"output_image_dir"`
	ComfyUIAddr    string `json:"comfy_ui_addr"`
}

var config Config

func initConfig() Config {
	file, err := os.ReadFile("./config.json")
	if err != nil {
		log.Println("Error reading JSON file:", err)
		return Config{}
	}
	var data Config
	err = json.Unmarshal(file, &data)
	if err != nil {
		log.Println("Error parsing JSON:", err)
		return Config{}
	}
	return data
}

func saveImage(imgBase64 string) {
	data, err := base64.StdEncoding.DecodeString(imgBase64)
	println("save image", imgBase64)
	if err != nil {
		log.Println("Error decoding image:", err)
		return
	}
	hasher := md5.New()
	hasher.Write(data)
	fileName := hex.EncodeToString(hasher.Sum(nil))
	println("save file", config.InputImageDir+fileName+".jpg")
	err = os.WriteFile(config.InputImageDir+fileName+".jpg", data, 0644)
}
func imageToBase64(imagePath string) string {
	imgFile, err := os.Open(imagePath)
	if err != nil {
		log.Fatalf("Failed to open image file: %v", err)
	}
	defer imgFile.Close()

	// 读取文件内容
	imgData, err := io.ReadAll(imgFile)
	if err != nil {
		log.Fatalf("Failed to read image file: %v", err)
	}

	// 将文件内容转换为base64编码
	imgBase64 := base64.StdEncoding.EncodeToString(imgData)

	return imgBase64
}
func queuePrompt(data []byte) {
	_, err := http.Post(fmt.Sprintf("%s/prompt", config.ComfyUIAddr), "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Println("Error sending prompt:", err)
		return
	}
}
func main() {
	config = initConfig()
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	log.Printf("connecting to %s", config.Addr)

	c, _, err := websocket.DefaultDialer.Dial(config.Addr, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
			}
			var receivedData routes.ToGPU
			err = json.Unmarshal(message, &receivedData)
			if err != nil {
				log.Println("unmarshal:", err)
			}
			if receivedData.Type == "prompt" {
				for _, imgBase64 := range receivedData.Images {
					saveImage(imgBase64)
				}
				postData := map[string]interface{}{
					"prompt": receivedData.Prompt,
				}
				jsonData, err := json.Marshal(postData)
				if err != nil {
					log.Println("marshal:", err)
				}
				queuePrompt(jsonData)
			} else if receivedData.Type == "image" {

			}
			//log.Printf("recv: %s", receivedData)
		}
	}()

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
