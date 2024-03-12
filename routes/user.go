package routes

import (
	"comfyui_service/db"
	"comfyui_service/model"
	"comfyui_service/utils"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"log"
	"net/http"
)

func getWxUser(code string) (db.WxLoginSession, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", utils.Config.APPID, utils.Config.APPSecret, code)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("get err:", err)
		return db.WxLoginSession{}, err
	}
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("read err:", err)
		return db.WxLoginSession{}, err
	}
	var result db.WxLoginSession
	err = json.Unmarshal(respBytes, &result)
	if err != nil {
		log.Println("unmarshal err:", err)
		return db.WxLoginSession{}, err
	}
	return result, nil
}

func Login(c *gin.Context) {
	wxCode := c.Query("code")
	wxUser, err := getWxUser(wxCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if wxUser.ErrCode != 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": wxUser.ErrMsg})
		return
	}
	var user db.User
	user, err = db.FindUser(wxUser)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			user, err = db.AddUser(wxUser)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
	}
	c.JSON(http.StatusOK, model.Response{Code: 0, Msg: "success", Data: user})
}
