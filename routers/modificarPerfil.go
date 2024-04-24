package routers

import (
	"context"
	"encoding/json"

	"github.com/LuisAli22/twitterGo/bd"
	"github.com/LuisAli22/twitterGo/models"
)

func ModificarPerfil(ctx context.Context, claim models.Claim) models.ResAPI {
	var r models.ResAPI
	r.Status = 400

	var t models.Usuario
	body := ctx.Value(models.Key("body")).(string)
	err := json.Unmarshal([]byte(body), &t)
	if err != nil {
		r.Message = "Datos Incorrectos" + err.Error()
	}
	status, err := bd.ModificoRegistro(t, claim.Id.Hex())
	if err != nil {
		r.Message = "Ocurrió un error al intentar modificar el registro " + err.Error()
		return r
	}

	if !status {
		r.Message = "No se ha logrado modificar el registro del usuario. "
		return r
	}

	r.Status = 200
	r.Message = "Modificación de perfil OK !"
	return r
}
