package bd

import (
	"context"

	"github.com/LuisAli22/twitterGo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InsertoTweet(t models.GraboTweet) (string, bool, error) {
	ctx := context.TODO()
	db := MongoCN.Database(DatabaseName)
	col := db.Collection("tweet")

	registro := bson.M{
		"userid":  t.UserID,
		"mensaje": t.Mensaje,
		"fecha":   t.Fecha,
	}
	result, err := col.InsertOne(ctx, registro)
	if err != nil {
		return "", false, err
	}
	//en result graba el objID que creo al insertar
	objID, _ := result.InsertedID.(primitive.ObjectID)
	return objID.String(), true, nil
}
