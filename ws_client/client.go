package ws_client

import (
	"comfyui_service/db"
	"comfyui_service/utils"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/url"
	"time"
)

type MessageData struct {
	PromptId         string `json:"prompt_id"`
	ExceptionMessage string `json:"exception_message"`
	Value            int    `json:"value"`
	Max              int    `json:"max"`
}

type Message struct {
	Type string      `json:"type"`
	Data MessageData `json:"data"`
}

func InitWs() {
	u := url.URL{Scheme: "ws", Host: utils.Config.ComfyHost, Path: "/ws", RawQuery: "clientId=" + utils.ComfyClientId}

	var c *websocket.Conn
	var err error
	var done chan struct{}

	connect := func() {
		log.Printf("connecting to %s", u.String())
		done = make(chan struct{})
		c, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			log.Println("dial:", err)
			return
		}

		go func() {
			defer close(done)
			for {
				var msg Message
				err := c.ReadJSON(&msg)
				if err != nil {
					log.Println("read:", err)
					return
				}
				switch msg.Type {
				case "execution_start":
					filter := bson.M{"job_id": msg.Data.PromptId}
					updater := bson.M{"$set": bson.M{"status": 1}}
					db.UpdateImageOne(filter, updater)
				case "progress":
					filter := bson.M{"job_id": msg.Data.PromptId}
					updater := bson.M{"$set": bson.M{"current": msg.Data.Value, "total": msg.Data.Max, "status": 1}}
					db.UpdateImageOne(filter, updater)
				case "executed":
					filter := bson.M{"job_id": msg.Data.PromptId}
					updater := bson.M{"$set": bson.M{"status": 2}}
					db.UpdateImageOne(filter, updater)
				case "execution_error":
					filter := bson.M{"job_id": msg.Data.PromptId}
					updater := bson.M{"$set": bson.M{"status": -1, "err_msg": msg.Data.ExceptionMessage}}
					db.UpdateImageOne(filter, updater)
				}
			}
		}()
	}

	connect()
	go func() {

		heartbeat := time.NewTicker(5 * time.Second) // send heartbeat every 5 seconds
		defer heartbeat.Stop()
		defer c.Close()

		for {
			select {
			case <-done:
				log.Println("connection lost, trying to reconnect")
				c.Close()
				time.Sleep(5 * time.Second)
				connect()
			case <-heartbeat.C:
				// Send heartbeat
				if c == nil {
					connect()
					continue
				}
				c.WriteMessage(websocket.PingMessage, nil)
			}
		}
	}()
}
