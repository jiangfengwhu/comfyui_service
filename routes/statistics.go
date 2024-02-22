package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ALive(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "我在"})
}
