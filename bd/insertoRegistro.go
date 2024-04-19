package bd

import (
	"context"

	"github.com/LuisAli22/twitterGo/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InsertoRegistro(u models.Usuario) (string, bool, error) {
	ctx := context.TODO()

	db := MongoCN.Database(DatabaseName)
	//Configuro con que colecion voy a trabajar. En este caso
	//voy a trabajar con la coleccion usuarios
	//las colecciones son como las tablas en sql
	// y Databasename es la base de datos de mongo
	col := db.Collection("usuarios")

	u.Password, _ = EncriptarPassword(u.Password)
	result, err := col.InsertOne(ctx, u)
	if err != nil {
		return "", false, err
	}

	ObjID, _ := result.InsertedID.(primitive.ObjectID)
	return ObjID.String(), true, nil
}
