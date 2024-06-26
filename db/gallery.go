package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ImageBase struct {
	Id         primitive.ObjectID `bson:"_id" json:"id"`
	Owner      primitive.ObjectID `bson:"owner" json:"owner"`
	Width      int                `bson:"width" json:"width"`
	Height     int                `bson:"height" json:"height"`
	JobId      string             `bson:"job_id" json:"job_id"`
	Public     bool               `bson:"public" json:"public"`
	Status     int                `bson:"status" json:"status"`
	TemplateId string             `bson:"template_id" json:"template_id"`
	ErrMsg     string             `bson:"err_msg" json:"err_msg"`
	Current    int                `bson:"current" json:"current"`
	Total      int                `bson:"total" json:"total"`
}

func AddImage(image ImageBase) error {
	_, err := db.gallery.InsertOne(context.TODO(), image)
	return err
}

func UpdateImageOne(filter interface{}, data interface{}) (*mongo.UpdateResult, error) {
	return db.gallery.UpdateOne(context.TODO(), filter, data)
}

func FindImage(filter interface{}, opts ...*options.FindOptions) ([]ImageBase, error) {
	var images []ImageBase
	cursor, err := db.gallery.Find(context.TODO(), filter, opts...)
	err = cursor.All(context.Background(), &images)
	return images, err
}

func DeleteImage(filter interface{}) (*mongo.DeleteResult, error) {
	return db.gallery.DeleteOne(context.TODO(), filter)
}
