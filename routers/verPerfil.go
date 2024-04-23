package routers

import (
	"encoding/json"
	"fmt"

	"github.com/LuisAli22/twitterGo/bd"
	"github.com/LuisAli22/twitterGo/models"
	"github.com/aws/aws-lambda-go/events"
)

func VerPerfil(request events.APIGatewayProxyRequest) models.ResAPI {
	var r models.ResAPI
	r.Status = 400
	fmt.Println("Entrè en VerPerfil")
	ID := request.QueryStringParameters["id"]
	if len(ID) < 1 {
		r.Message = "El paràmetro ID es obligatorio"
		return r
	}

	perfil, err := bd.BuscoPerfil(ID)
	if err != nil {
		r.Message = "Ocurriò un error al intentar buscar el registro" + err.Error()
		return r
	}
	respJson, err := json.Marshal(perfil)
	if err != nil {
		r.Status = 500
		r.Message = "Error al formatoear los datos de los usuarios como JSON " + err.Error()
		return r
	}

	r.Status = 200
	r.Message = string(respJson)
	return r
}
