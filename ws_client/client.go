package ws_client

import (
	"comfyui_service/db"
	"comfyui_service/utils"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/url"
)

type MessageData struct {
	PromptId         string `json:"prompt_id"`
	ExceptionMessage string `json:"exception_message"`
}

type Message struct {
	Type string      `json:"type"`
	Data MessageData `json:"data"`
}

func InitWs() {
	u := url.URL{Scheme: "ws", Host: utils.Config.ComfyHost, Path: "/ws", RawQuery: "clientId=" + utils.ComfyClientId}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	go func() {
		for {
			var msg Message
			err := c.ReadJSON(&msg)
			if err != nil {
				log.Println("read:", err)
				return
			}
			if msg.Type == "execution_start" {
				filter := bson.M{"job_id": msg.Data.PromptId}
				updater := bson.M{"$set": bson.M{"status": 1}}
				_, err := db.UpdateImageOne(filter, updater)
				if err != nil {
					log.Println("update err:", err.Error())
				}
			}
			if msg.Type == "executed" {
				filter := bson.M{"job_id": msg.Data.PromptId}
				updater := bson.M{"$set": bson.M{"status": 2}}
				_, err := db.UpdateImageOne(filter, updater)
				if err != nil {
					log.Println("update err:", err.Error())
				}
			}
			if msg.Type == "execution_error" {
				filter := bson.M{"job_id": msg.Data.PromptId}
				updater := bson.M{"$set": bson.M{"status": -1, "err_msg": msg.Data.ExceptionMessage}}
				_, err := db.UpdateImageOne(filter, updater)
				if err != nil {
					log.Println("update err:", err.Error())
				}
			}
		}
	}()
}
