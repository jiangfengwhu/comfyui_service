package routes

import (
	"comfyui_service/db"
	"comfyui_service/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

func HomeGallery(c *gin.Context) {
	filter := bson.M{"public": true}
	data, err := db.FindImage(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.Response{Code: 0, Msg: "success", Data: data})
}
