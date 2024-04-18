package models

/*
Este struct lleva los campos que ya incluimos
en nuestro secret en aws.
Tenemos que ponerlo con mayuscula a los campos
porque necesitamos que sean publicos y usarlos
en otro paquete
*/
type Secret struct {
	Host     string `json: "host"`
	Username string `json: "username"`
	Password string `json: "password"`
	JWTSign  string `json: "jwtsign"`
	Database string `json: "database"`
}
