package bd

import (
	"context"

	"github.com/LuisAli22/twitterGo/models"
	"go.mongodb.org/mongo-driver/bson"
)

func ChequeoYaExisteUsuario(email string) (models.Usuario, bool, string) {
	ctx := context.TODO()

	db := MongoCN.Database(DatabaseName)
	//Configuro con que colecion voy a trabajar. En este caso
	//voy a trabajar con la coleccion usuarios
	//las colecciones son como las tablas en sql
	// y Databasename es la base de datos de mongo
	col := db.Collection("usuarios")
	//Tenemos que buscar dentro de mongo si existe el usuario
	//Para eso mongo usa condiciones, pero las condiciones tienen
	//que tener un formato especial.
	//bson.M es una interfaz de especifica de bson. Va a tener un formato
	//clave-valor. Le vamos a decir que tiene que ser filtrado por el mail.
	//Es como el Where de sql
	condition := bson.M{"email": email}
	var resultado models.Usuario
	//Me interesa obtener el primero que encuentre
	//FindOne recibe el contexto y la condicion y despues lo que obtiene lo
	//codifica en la variable resultad  (por eso le paso el puntero a resultado)
	err := col.FindOne(ctx, condition).Decode(&resultado)
	//Tomo el valor hexadecimal de ID y lo convierto a string
	ID := resultado.ID.Hex()
	if err != nil {
		return resultado, false, ID
	}
	return resultado, true, ID
}
