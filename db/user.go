package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type WxLoginSession struct {
	OpenId     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionId    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

type UserWx struct {
	OpenId     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionId    string `json:"unionid"`
}

type User struct {
	NickName string             `bson:"nick_name" json:"nick_name"`
	Avatar   string             `bson:"avatar" json:"avatar"`
	Id       primitive.ObjectID `bson:"_id" json:"id"`
	Tickets  int                `bson:"tickets" json:"tickets"`
}

func FindUser(user WxLoginSession) (User, error) {
	var userResp User
	err := db.user.FindOne(context.TODO(), bson.D{{"openid", user.OpenId}}).Decode(&userResp)
	return userResp, err
}

func AddUser(user WxLoginSession) (User, error) {
	userResp := User{NickName: "小土豆", Avatar: "", Tickets: 10}
	result, err := db.user.InsertOne(context.TODO(), bson.D{{"nick_name", userResp.NickName}, {"avatar", userResp.Avatar}, {"tickets", userResp.Tickets}, {"openid", user.OpenId}, {"session_key", user.SessionKey}, {"unionid", user.UnionId}})
	if err != nil {
		return User{}, err
	}
	userResp.Id = result.InsertedID.(primitive.ObjectID)
	return userResp, nil
}

func UpdateUserOne(filter interface{}, data interface{}) (*mongo.UpdateResult, error) {
	return db.user.UpdateOne(context.TODO(), filter, data, options.Update().SetUpsert(true))
}
