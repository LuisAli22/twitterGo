package routers

import (
	"context"

	"github.com/LuisAli22/twitterGo/bd"
	"github.com/LuisAli22/twitterGo/models"
	"github.com/aws/aws-lambda-go/events"
)

func AltaRelacion(ctx context.Context, request events.APIGatewayProxyRequest, claim models.Claim) models.ResAPI {
	var r models.ResAPI
	r.Status = 400
	ID := request.QueryStringParameters["id"]
	if len(ID) < 1 {
		r.Message = "El paràmetro ID es obligatorio"
		return r
	}

	var t models.Relacion
	t.UsuarioID = claim.Id.Hex()
	t.UsuarioRelacionID = ID

	status, err := bd.InsertoRelacion(t)
	if err != nil {
		r.Message = "Ocurrio un error al intentar insertar relación " + err.Error()
		return r
	}
	if !status {
		r.Message = "No se ha logrado insertar la relación "
		return r
	}

	r.Status = 200
	r.Message = "Alta de relacion Ok!"
	return r

}
