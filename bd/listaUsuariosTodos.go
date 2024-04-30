package bd

import (
	"context"
	"fmt"

	"github.com/LuisAli22/twitterGo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func LeoUsuariosTodos(ID string, page int64, search string, tipo string) ([]*models.Usuario, bool) {
	ctx := context.TODO()
	db := MongoCN.Database(DatabaseName)
	col := db.Collection("usuarios")

	var results []*models.Usuario
	opciones := options.Find()
	opciones.SetLimit(20)                               //Va enviar paginado de 20 tweet cada uno
	opciones.SetSort(bson.D{{Key: "fecha", Value: -1}}) //lo va a ordenar por fecha
	//y como value es -1 lo ordena de forma descendente
	//poniendo Value: 1 lo ordena ascendente
	opciones.SetSkip((page - 1) * 20) //es el OFFSET de SQL para decir por ejemplo,
	//que me devuelvea 20 registros que indico en SetLimit pero salteando los primeros
	//k registros, donde k es el parametro que recibe el Skip
	//Ahor hacemos realmente el find
	query := bson.M{
		"nombre": bson.M{"$regex": `(?i)` + search},
	}

	cur, err := col.Find(ctx, query, opciones)
	if err != nil {
		return results, false
	}

	var incluir bool

	for cur.Next(ctx) {
		var s models.Usuario
		err := cur.Decode(&s)
		if err != nil {
			fmt.Println("DEcode = " + err.Error())
			return results, false
		}
		var r models.Relacion
		r.UsuarioID = ID
		r.UsuarioRelacionID = s.ID.Hex()

		incluir = false
		encontrado := ConsultoRelacion(r)
		if tipo == "new" && !encontrado {
			incluir = true
		}
		if tipo == "follow" && encontrado {
			incluir = true
		}

		if r.UsuarioRelacionID == ID {
			incluir = false
		}
		if incluir {
			s.Password = "" //no quiero enviar el password de cada uno. Por eso le pongo cadena vacia
			results = append(results, &s)
		}
	}
	err = cur.Err()
	if err != nil {
		fmt.Println("cur.Err() = " + err.Error())
		return results, false
	}
	cur.Close(ctx)
	return results, true
}
