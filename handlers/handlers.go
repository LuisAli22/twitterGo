package handlers

import (
	"context"
	"fmt"

	"github.com/LuisAli22/twitterGo/models"
	"github.com/aws/aws-lambda-go/events"
)

func Manejadores(ctx context.Context, request events.APIGatewayProxyRequest) models.ResAPI {
	//Quiero enviar a Cloudwatch un mensaje que dice: Voy a procesar el "path" con tal metodo
	fmt.Println("Voy a procesar " + ctx.Value(models.Key("path")).(string) + " > " + ctx.Value(models.Key("method")).(string))
	var r models.ResAPI
	//Le da un valor inicial de 400 porque son muchas mas las veces que va a dar error, que va a faltar
	//que simplemente cuando este todo ok
	r.Status = 400

	switch ctx.Value(models.Key("method")).(string) {
	case "POST":
		switch ctx.Value(models.Key("path")).(string) {

		}
	case "GET":
		switch ctx.Value(models.Key("path")).(string) {

		}
	case "PUT":
		switch ctx.Value(models.Key("path")).(string) {

		}
	case "DELETE":
		switch ctx.Value(models.Key("path")).(string) {

		}
	}
	r.Message = "Method Invalid"
	return r
}
