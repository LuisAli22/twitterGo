package routers

import (
	"github.com/LuisAli22/twitterGo/bd"
	"github.com/LuisAli22/twitterGo/models"
	"github.com/aws/aws-lambda-go/events"
)

func BajaRelacion(request events.APIGatewayProxyRequest, claim models.Claim) models.ResAPI {
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

	status, err := bd.BorroRelacion(t)
	if err != nil {
		r.Message = "Ocurriò un error al intentar borrar relaciòn"
		return r
	}
	if !status {
		r.Message = "No se ha logrado borrar relaciòn"
		return r
	}
	r.Status = 200
	r.Message = "Baja Relaciòn OK!"
	return r
}
