package main

import (
	"bytes"
	"comfyui_service/gpu_host"
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
	"sync"
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
	mutex := sync.Mutex{}
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	var c *websocket.Conn
	var err error
	var done chan struct{}

	connect := func() {
		log.Printf("connecting to %s", config.Addr)
		done = make(chan struct{})
		c, _, err = websocket.DefaultDialer.Dial(config.Addr, nil)
		if err != nil {
			log.Println("dial:", err)
			return
		}

		go func() {
			defer close(done)
			for {
				var receivedData gpu_host.WSMessage
				err := c.ReadJSON(&receivedData)
				if err != nil {
					log.Println("error msg:", err)
					return
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
						return
					}
					queuePrompt(jsonData)
				} else if receivedData.Type == "image" {

				} else if receivedData.Type == "alive" {
					mutex.Lock()
					err := c.WriteJSON(gpu_host.WSMessage{Type: "alive", Id: receivedData.Id})
					if err != nil {
						log.Println("write alive err:", err)
						return
					}
					mutex.Unlock()
				}
			}
		}()
	}

	connect()
	heartbeat := time.NewTicker(5 * time.Second) // send heartbeat every 5 seconds
	defer heartbeat.Stop()

	for {
		select {
		case <-done:
			log.Println("connection lost, trying to reconnect")
			c.Close()
			time.Sleep(5 * time.Second)
			connect()
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
		case <-heartbeat.C:
			// Send heartbeat
			if c == nil {
				connect()
				continue
			}
			err := c.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				log.Println("write heartbeat:", err)
				c.Close()
				time.Sleep(5 * time.Second)
				connect()
			}
		}
	}
}
