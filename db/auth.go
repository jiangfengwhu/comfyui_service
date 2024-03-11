package db

import (
	"comfyui_service/model"
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type UserModel struct {
	SessionKey string `json:"session_key" bson:"session_key"`
}

func AuthRequired(c *gin.Context) {
	token := c.GetHeader("Bearer")
	objectId, err := primitive.ObjectIDFromHex(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.Response{Code: 403, Msg: "未登录"})
		return
	}
	var user UserModel
	err = db.user.FindOne(context.TODO(), bson.M{"_id": objectId}).Decode(&user)
	c.Set("user", user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.Response{Code: 403, Msg: "未登录"})
		return
	}
	c.Next()
}
