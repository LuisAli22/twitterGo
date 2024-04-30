package routers

import (
	"encoding/json"
	"strconv"

	"github.com/LuisAli22/twitterGo/bd"
	"github.com/LuisAli22/twitterGo/models"
	"github.com/aws/aws-lambda-go/events"
)

func ListaUsuarios(request events.APIGatewayProxyRequest, claim models.Claim) models.ResAPI {
	var r models.ResAPI
	r.Status = 400

	page := request.QueryStringParameters["page"]
	typeUser := request.QueryStringParameters["type"]
	search := request.QueryStringParameters["search"]
	IDUsuario := claim.Id.Hex()

	if len(page) == 0 {
		page = "1"
	}

	pageTemp, err := strconv.Atoi(page)
	if err != nil {
		r.Message = "Debe enviar el paràmetro 'page' como entero mayor a cero " + err.Error()
		return r
	}
	//mongo pide que pageTemp sea un int64
	usuarios, status := bd.LeoUsuariosTodos(IDUsuario, int64(pageTemp), search, typeUser)
	if !status {
		r.Message = "Error al leer los susuarios"
		return r
	}

	respJson, err := json.Marshal(usuarios)
	if err != nil {
		r.Status = 500
		r.Message = "Error al formatear los datos de los usuarios en JSON"
	}

	r.Status = 200
	r.Message = string(respJson)
	return r
}
