package routes

import (
	"comfyui_service/db"
	"comfyui_service/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strconv"
)

type ImageDelReq struct {
	Id primitive.ObjectID `json:"id" binding:"required"`
}

func HomeGallery(c *gin.Context) {
	filter := bson.M{"public": true}
	data, err := db.FindImage(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.Response{Code: 0, Msg: "success", Data: data})
}

func DeleteImage(c *gin.Context) {
	var req ImageDelReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code: -1,
			Msg:  err.Error(),
		})
		return
	}
	user := c.MustGet("user").(db.User)
	filter := bson.M{"_id": req.Id, "owner": user.Id}
	result, err := db.DeleteImage(filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code: -1,
			Msg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, model.Response{Code: 0, Msg: strconv.FormatInt(result.DeletedCount, 10)})
}
