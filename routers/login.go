package routers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/LuisAli22/twitterGo/bd"
	"github.com/LuisAli22/twitterGo/jwt"
	"github.com/LuisAli22/twitterGo/models"
	"github.com/aws/aws-lambda-go/events"
)

func Login(ctx context.Context) models.ResAPI {
	var t models.Usuario
	var r models.ResAPI
	r.Status = 400

	body := ctx.Value(models.Key("body")).(string)
	err := json.Unmarshal([]byte(body), &t)
	if err != nil {
		r.Message = "Usuario y/o contraseña Invàlidos" + err.Error()
		return r
	}
	if len(t.Email) == 0 {
		r.Message = "El email del usuario es requerido"
		return r
	}

	userData, existe := bd.IntentoLogin(t.Email, t.Password)
	if !existe {
		r.Message = "Usuario y/o contraseña Inválidos"
		return r
	}
	jwtKey, err := jwt.GeneroJWT(ctx, userData)
	if err != nil {
		r.Message = "Ocurriò un error al intentar generar el token correspondiente" + err.Error()
		return r

	}
	resp := models.RespuestaLogin{
		Token: jwtKey,
	}

	token, err2 := json.Marshal(resp)
	if err2 != nil {
		r.Message = "Ocurriò un error al intentar formatear el token a JSON" + err2.Error()
		return r
	}

	//Ademas del token, vamos a devolver una cookie para que el usuario tenga el token grabado en una cookie
	//del sistema y cuando vuelva a loguearse no tenga que ir a buscar otro token
	cookie := &http.Cookie{
		Name:    "token",
		Value:   jwtKey,
		Expires: time.Now().Add(24 * time.Hour),
	}

	cookieStr := cookie.String()
	res := &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(token),
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
			"Set-Cookie":                  cookieStr,
		},
	}
	r.Status = 200
	r.Message = string(token)
	r.CustomResp = res
	return r
}
