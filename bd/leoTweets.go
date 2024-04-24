package bd

import (
	"context"

	"github.com/LuisAli22/twitterGo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Devuelve un slice de n putnero a models.DevuelvoTweets y un boolean
func LeoTweets(ID string, pagina int64) ([]*models.DevuelvoTweets, bool) {
	ctx := context.TODO()
	db := MongoCN.Database(DatabaseName)
	col := db.Collection("tweet")
	//slice vacio para resultados
	var resultados []*models.DevuelvoTweets
	//Tengo un ID de usuario y quiero leer los tweets de ese ID de usuario
	//
	condicion := bson.M{
		"userid": ID,
	}
	//Mongo me permite configurar opciones. Dentro de un llamado, cuando
	//hago un find() me permite incluir opciones que en  SQL vendrian a ser
	//partes del WHERE, tiene partes que serian el GROUP BY, el ORDER BY
	//En este caso todo se maneja bajo opciones
	//A continucacion le dio que las opciones que voy a configurar son para el FIND()

	opciones := options.Find()
	opciones.SetLimit(20)                               //Va enviar paginado de 20 tweet cada uno
	opciones.SetSort(bson.D{{Key: "fecha", Value: -1}}) //lo va a ordenar por fecha
	//y como value es -1 lo ordena de forma descendente
	//poniendo Value: 1 lo ordena ascendente
	opciones.SetSkip((pagina - 1) * 20) //es el OFFSET de SQL para decir por ejemplo,
	//que me devuelvea 20 registros que indico en SetLimit pero salteando los primeros
	//k registros, donde k es el parametro que recibe el Skip
	//Ahor hacemos realmente el find
	cursor, err := col.Find(ctx, condicion, opciones)
	if err != nil {
		return resultados, false
	}
	//U cursor es un conjunto de registros. Tengo que iterar el cursor para ir
	//armando el slice a devolver
	for cursor.Next(ctx) {
		var registro models.DevuelvoTweets
		err := cursor.Decode(&registro)
		if err != nil {
			return resultados, false
		}
		resultados = append(resultados, &registro)
	}
	return resultados, true
}
