package models

/*
Esto es porque el front pide un true o un false. Estoy relacionado
con este ID de usuario? (Si o no)
*/

type RespuestaConsultaRelacion struct {
	Status bool `json:"status"`
}
