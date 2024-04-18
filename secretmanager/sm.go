package secretmanager

import (
	"encoding/json"
	"fmt"

	"github.com/LuisAli22/twitterGo/awsgo"
	"github.com/LuisAli22/twitterGo/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

/*
Yo recibia por variable de entorno el secret name. Eso que recibimos
o vamos a usar como parametro para esta funcion.
Nos devuelve la estrucutura con los campos Host, Username, Password, etc
y en error en caso de qe suceda.
*/
func GetSecret(secretName string) (models.Secret, error) {
	var datosSecret models.Secret
	fmt.Println("> Pido Secreto" + secretName)
	//Inicializamos secret manager
	svc := secretsmanager.NewFromConfig(awsgo.Cfg)
	//Capturo la clave (conjunto clave-valor)
	clave, err := svc.GetSecretValue(awsgo.Ctx, &secretsmanager.GetSecretValueInput{
		// Amazon tiene algunas funciones que manejan los string de forma diferente a los string de Go
		//Tienen un formato diferente, por eso no puedo pasarle el literal. Lo tengo que convertir
		// en aws.String(). Es una funciòn de amazon para convertirlo al formato de ellos
		SecretId: aws.String(secretName),
	})
	if err != nil {
		fmt.Println(err.Error())
		return datosSecret, err
	}
	//convierto el secret decodificado en la estructura json
	json.Unmarshal([]byte(*clave.SecretString), &datosSecret)
	fmt.Println("> Lectura de Secret OK " + secretName)
	return datosSecret, nil
}
