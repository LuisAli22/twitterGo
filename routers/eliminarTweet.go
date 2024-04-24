package routers

import (
	"github.com/LuisAli22/twitterGo/bd"
	"github.com/LuisAli22/twitterGo/models"
	"github.com/aws/aws-lambda-go/events"
)

func EliminarTweet(request events.APIGatewayProxyRequest, claim models.Claim) models.ResAPI {
	var r models.ResAPI
	r.Status = 400

	ID := request.QueryStringParameters["id"]
	if len(ID) < 1 {
		r.Message = "El parámetro ID es obligatorio"
		return r
	}

	err := bd.BorroTweet(ID, claim.Id.Hex())
	if err != nil {
		r.Message = "Ocurrió un error al intentar borrar el tweet " + err.Error()
		return r
	}
	r.Message = "Eliminar Tweet OK !"
	r.Status = 200
	return r
}
