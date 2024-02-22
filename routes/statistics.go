package routes

import (
	"comfyui_service/gpu_host"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"sync"
)

func ALive(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "我在"})
}

var lock = sync.Mutex{}

func GPUALive(c *gin.Context) {
	var receivedMsg gpu_host.WSMessage
	id := uuid.New().String()
	pingMsg := gpu_host.WSMessage{Type: "alive", Id: id}
	msgChan := make(chan string)
	go func() {
		for {
			err := gpu_host.GpuConnection.Conn.ReadJSON(&receivedMsg)
			if err != nil {
				msgChan <- "GPU连接错误"
				return
			}
			if receivedMsg.Id == id {
				msgChan <- "我在"
				return
			}
		}
	}()
	err := gpu_host.GpuConnection.SendMessage(pingMsg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "发送失败"})
		return
	}
	msg := <-msgChan
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": msg})
}
