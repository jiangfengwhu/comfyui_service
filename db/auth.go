package db

import (
	"comfyui_service/model"
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

func AuthRequired(c *gin.Context) {
	token := c.GetHeader("Bearer")
	var user User
	err := db.user.FindOne(context.TODO(), bson.M{"openid": token}).Decode(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, model.Response{Code: 401, Msg: "未登录"})
		return
	}
	c.Set("user", user)
	c.Next()
}
