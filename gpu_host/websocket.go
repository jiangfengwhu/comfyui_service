package gpu_host

import (
	"comfyui_service/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"sync"
)

type WSMessage struct {
	Type    string                    `json:"type"`
	Prompt  map[string]utils.BaseNode `json:"prompt,omitempty"`
	Images  map[string]string         `json:"images,omitempty"`
	Id      string                    `json:"id,omitempty"`
	Content string                    `json:"content,omitempty"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type GpuConnectionType struct {
	Conn       *websocket.Conn
	ReadMutex  sync.Mutex
	WriteMutex sync.Mutex
}

func (c *GpuConnectionType) SendMessage(v WSMessage) error {
	c.WriteMutex.Lock()
	defer c.WriteMutex.Unlock()
	return c.Conn.WriteJSON(v)
}

//	func (c *GpuConnectionType) ReceiveMessage(Id string, v *WSMessage) error {
//		c.ReadMutex.Lock()
//		defer c.ReadMutex.Unlock()
//		deadline := time.Now().Add(3 * time.Second)
//
//		for {
//			// Read the next JSON-encoded message from the connection
//			err := c.Conn.ReadJSON(v)
//			if err != nil {
//				return err
//			}
//			println("received message", v.Id)
//			// Check if the Id of the message matches the input Id
//			if v.Id == Id {
//				return nil
//			}
//			// If the deadline has passed, return nil
//			if time.Now().After(deadline) {
//				return nil
//			}
//		}
//	}
func (c *GpuConnectionType) setConn(conn *websocket.Conn) {
	c.Conn = conn
}
func (c *GpuConnectionType) isConnected() bool {
	return c.Conn != nil
}

var GpuConnection = GpuConnectionType{
	Conn:       nil,
	ReadMutex:  sync.Mutex{},
	WriteMutex: sync.Mutex{},
}

func LinkRemoteGPU(c *gin.Context) {
	if GpuConnection.isConnected() {
		c.JSON(400, gin.H{"msg": "已经有连接了"})
		return
	}
	conn, _ := upgrader.Upgrade(c.Writer, c.Request, nil)
	GpuConnection.setConn(conn)
}
