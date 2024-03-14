package ws_client

import (
	"comfyui_service/db"
	"comfyui_service/utils"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"
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
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

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
					c.Close()
				}
				return
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
