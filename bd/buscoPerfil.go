package bd

import (
	"context"

	"github.com/LuisAli22/twitterGo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func BuscoPerfil(ID string) (models.Usuario, error) {
	ctx := context.TODO()
	db := MongoCN.Database(DatabaseName)
	col := db.Collection("usuarios")
	var perfil models.Usuario
	//para buscr en Mongo tenemos que convertir el ID de string a primitive
	objID, _ := primitive.ObjectIDFromHex(ID)

	condicion := bson.M{
		"_id": objID,
	}
	err := col.FindOne(ctx, condicion).Decode(&perfil)
	//no me interesa mandar el password al front
	//Le pongo cadena vacia asi despues el moit empty lo detecta y no lo
	//pone en el json
	perfil.Password = ""
	if err != nil {
		return perfil, err
	}
	return perfil, nil
}
