package bd

import (
	"context"
	"fmt"

	"github.com/LuisAli22/twitterGo/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// En este variable vamos a guardar en esta variable la conexion abierta
// hacia mongo
var MongoCN *mongo.Client
var DatabaseName string

func ConectarBD(ctx context.Context) error {
	//A context.Value() recibe el tipo models.Key que definimos en tyoes.go
	//despues con .(string) lo transforma a string.
	//la variable user pasa a ser un string comun de GO
	user := ctx.Value(models.Key("user")).(string)
	passwd := ctx.Value(models.Key("password")).(string)
	host := ctx.Value(models.Key("host")).(string)
	//Eñ fmt no solo sirve para loguear sino tambien para formatear texto
	//Sprintf -> La S inicial indica que devuelve un string (diferencia entre Printf y Sprintf)
	connStr := fmt.Sprintf("mongodb+srv://%s:%s@%s/?retryWrites=true&w=majority", user, passwd, host)
	var clientOptions = options.Client().ApplyURI(connStr)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	//Una vez que me conecto hago un ping para ver si la conexion quedo abierta o
	//tuvo algun error que no pude capturar
	err = client.Ping(ctx, nil)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Conexión exitosa con la BD")
	//Mongo client es el client cuando me conecto
	MongoCN = client
	DatabaseName = ctx.Value(models.Key("database")).(string)
	return nil
}

func BaseConectada() bool {
	err := MongoCN.Ping(context.TODO(), nil)
	return err == nil
}
