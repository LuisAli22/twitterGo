package bd

import (
	"context"

	"github.com/LuisAli22/twitterGo/models"
	"go.mongodb.org/mongo-driver/bson"
)

func LeoTweetsSeguidores(ID string, pagina int) ([]models.DevuelvoTweetsSeguidores, bool) {
	ctx := context.TODO()
	db := MongoCN.Database(DatabaseName)
	col := db.Collection("relacion")
	skip := (pagina - 1) * 20

	/*
		Crea un slice de bson.M porque van a ser varias condiciones
		Se crea con un size igual a 0 (cero)para despues ir agregando con append
	*/
	condiciones := make([]bson.M, 0)
	/*
		A continuacion se hace un Join en Mongo.
		No es le fuerte de mongo hacer reelaciones entre colecciones.
		En Mongo es muy complicado hacer un join y no e sperformante tampoco
	*/
	//Match es la clave de bùsqueda y como valor le coloca un bson.M con el ID
	condiciones = append(condiciones, bson.M{"$match": bson.M{"usuarioid": ID}})
	//Nueva condicion: Para unir dos colecciones en Mongo se usa la instruccion $lookup
	//Para hacer el UNION como en SQl hay que decirle cuales son los campos para unir
	//Entonces $lookup tiene un bson.M con todas las condiciones
	//Se hace un append para conservar lo que ya estaba en colecciones
	//la primera es FROM para indicar desde que coleccion
	//despues, tenemos que decirle cual va a ser el localfield (local en este caso
	//va a ser la tabla relacion). Entonces el localfield es el campo usuariorelacionid.
	//PAra la Foreign key tenemos que indicar quien es el foreignField, que en este caso
	//va a ser el campo userid (de la tabla tweets)
	// con AS se pone un alias. Indicamos que el alias es tweet
	condiciones = append(condiciones, bson.M{
		"$lookup": bson.M{
			"from":         "tweet",
			"localField":   "usuariorelacionid",
			"foreignField": "userid",
			"as":           "tweet",
		},
	})
	//Agregamos otras condiciones
	//
	condiciones = append(condiciones, bson.M{"$unwind": "$tweet"})
	condiciones = append(condiciones, bson.M{"$sort": bson.M{"tweet.fecha": -1}})
	condiciones = append(condiciones, bson.M{"$skip": skip})
	//limit es cuantos tweet voy a mostrar por pagina
	condiciones = append(condiciones, bson.M{"$limit": 20})
	var result []models.DevuelvoTweetsSeguidores
	/*El Aggregate permite incorporar a la coleccion
	todoslos filtros y condiciones que definimos para que el cursor
	venga con todos los filtros aplicados
	*/
	cursor, err := col.Aggregate(ctx, condiciones)
	if err != nil {
		return result, false
	}

	//Cuando devolvemos el cursor se convierte en un objeto de Mongo
	//A continuacion usamos una instruccion de cursor
	//Vamos a devolver todos los resultados con cursor.All()
	//Esto toma todo del cursor y lo pone en result que es nuestro slice de
	// models.DevuelvoTweetsSeguidores
	err = cursor.All(ctx, &result)
	if err != nil {
		return result, false
	}
	return result, true

}
