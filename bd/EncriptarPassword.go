package bd

import "golang.org/x/crypto/bcrypt"

func EncriptarPassword(pass string) (string, error) {
	//bcrypt encripta la password. Es encirptacion de una sola via. La encripta pero
	//no se desencripta nunca. Es como el md5. A lo sumo cuando el usuario se loguea
	// y llega el password del usuario, la encripta n veces (n es el costo, en este caso 8) y
	// eso se compara con lo que esta guardado en la base de datos
	//El costo es la cantidad de vueltas que va a realizar el modulo de encriptacion para
	//encriptar el password. Mientras mas alto es el costo, mas segura es la contraseña,
	//pero usa mas recursos, demora mas en encriptar
	costo := 8
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), costo)
	if err != nil {
		return err.Error(), err
	}
	return string(bytes), nil
}
