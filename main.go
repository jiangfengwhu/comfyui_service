package main

import (
	"comfyui_service/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	//"github.com/tidwall/sjson"
)

func main() {
	r := gin.Default()
	var data = utils.ReadWorkflowFile()
	//fmt.Println(gjson.Get(data, `..#.class_type`))
	gjson.AddModifier("case", func(json, arg string) string {
		if arg == "upper" {
			return strings.ToUpper(json)
		}
		if arg == "lower" {
			return strings.ToLower(json)
		}
		return json
	})
	println(gjson.Get(data, `3|@case:upper`).String())
	CreateKSampler := utils.CreateKSampler()
	CreateKSampler.UpdateSeed()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
