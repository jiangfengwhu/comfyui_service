package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
var gpuConnection *websocket.Conn

func LinkRemoteGPU(c *gin.Context) {
	conn, _ := upgrader.Upgrade(c.Writer, c.Request, nil)
	gpuConnection = conn
	for {
		// Read message from client
		messageType, p, _ := conn.ReadMessage()
		// Echo message back to client
		conn.WriteMessage(messageType, p)
	}
}
