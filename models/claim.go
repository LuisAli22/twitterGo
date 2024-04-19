package models

//el claim se refiere al payload
//bson: es un json encriptado
import (
	jwt "github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Claim struct {
	Email string `json:"email"`
	//primitive.ObjectId es un tipo de dato especifico existente para un Id
	//Mongo pone al id como _id (con uion bajo) y tiene que ser un bson
	//A su vez le decimos que el formato en json (cuando este ya descompactado)
	// es _id y le agregamos omitempty para que cuando ecuente este dato vacio, lo omita
	//No lo va a incluir en la estructura del json
	Id primitive.ObjectID `bson:"_id" json:"_id,omitempty"`
	/*
		Le digo que el resto del json lo complete con el registered claim
	*/
	jwt.RegisteredClaims
}
