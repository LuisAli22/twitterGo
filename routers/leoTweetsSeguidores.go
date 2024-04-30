package routers

import (
	"encoding/json"
	"strconv"

	"github.com/LuisAli22/twitterGo/bd"
	"github.com/LuisAli22/twitterGo/models"
	"github.com/aws/aws-lambda-go/events"
)

func LeoTweetsSeguidores(request events.APIGatewayProxyRequest, claim models.Claim) models.ResAPI {
	var r models.ResAPI
	r.Status = 400
	IDUsuario := claim.Id.Hex()
	//Los tweets van paginados
	pagina := request.QueryStringParameters["pagina"]
	if len(pagina) < 1 {
		pagina = "1"
	}
	pag, err := strconv.Atoi(pagina)
	if err != nil {
		r.Message = "Debe enviar el paràmetro Pàgina como un valor mayor a 0"
		return r
	}

	tweets, correcto := bd.LeoTweetsSeguidores(IDUsuario, pag)
	if !correcto {
		r.Message = "Error al leer los tweets"
		return r
	}
	respJson, err := json.Marshal(tweets)
	if err != nil {
		r.Status = 500
		r.Message = "Error al formatear los datos de tweets de los eguidores"
	}

	r.Status = 200
	r.Message = string(respJson)
	return r

}
