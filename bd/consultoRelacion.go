package bd

import (
	"context"

	"github.com/LuisAli22/twitterGo/models"
	"go.mongodb.org/mongo-driver/bson"
)

func ConsultoRelacion(t models.Relacion) bool {
	ctx := context.TODO()
	db := MongoCN.Database(DatabaseName)
	col := db.Collection("relacion")

	condicion := bson.M{
		"usuarioid":         t.UsuarioID,
		"usuariorelacionid": t.UsuarioRelacionID,
	}

	var resultado models.Relacion
	err := col.FindOne(ctx, condicion).Decode(&resultado)
	//Cuando hace un findOne y no encuentra nada, Mongo devuelve un error
	//pero segun el que dicta el curso es algo que esta mal porque no
	//encontrar nada no es un error
	return err == nil
}
