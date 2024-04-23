package main

import (
	"context"
	"os"
	"strings"

	"github.com/LuisAli22/twitterGo/awsgo"
	"github.com/LuisAli22/twitterGo/bd"
	"github.com/LuisAli22/twitterGo/handlers"
	"github.com/LuisAli22/twitterGo/models"
	"github.com/LuisAli22/twitterGo/secretmanager"
	"github.com/aws/aws-lambda-go/events"
	lambda "github.com/aws/aws-lambda-go/lambda"
)

func main() {
	//La llamada y el comienzo de la lambda
	lambda.Start(HandleRequest)
}

// Cuando tengo una API REST en Amazon y llamo al API Gateway desde la aplicacion
// o desde Postman, la API gateway se va a conectar con mi lambda y le va a mandar
// un objeto de tipo APIGatewayProxyRequest
func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	//Aca va la logica de mi desarrollo.
	//Todo lo que pong aca va a ser lo que ejecute cuando entra en la lambda
	var res *events.APIGatewayProxyResponse
	awsgo.InicializoAWS()
	if !ValidoParametros() {
		res = &events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Error en las variables de entorno. Deben incluir 'SecretName', 'BucketName', 'UrlPrefix'",
			Headers: map[string]string{
				"Content-Type": "Application/Json",
			},
		}
		return res, nil
	}
	SecretModel, err := secretmanager.GetSecret(os.Getenv("SecretName"))
	if err != nil {
		res = &events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Error en la lectura de Secret " + err.Error(),
			Headers: map[string]string{
				"Content-Type": "Application/Json",
			},
		}
		return res, nil
	}
	//Cuando encuentra el prefijo que viene por parametro, lo reemplaza por la
	//cadena vacia
	path := strings.Replace(request.PathParameters["twittergo"], os.Getenv("UrlPrefix"), "", -1)
	// A continuacion colocamos variables en el context
	// El segundo parametro de WithValue pide un key que puede ser
	// de cualquier tipo. Podrìa poner un string, pero GO (y cualquier lenguaje)
	// es no utlizar explicitamente un tipo string para una clave. Tenemos que
	//crearnos un tipo de dato propio nuestro.
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("path"), path)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("method"), request.HTTPMethod)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("user"), SecretModel.Username)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("password"), SecretModel.Password)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("host"), SecretModel.Host)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("database"), SecretModel.Database)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("jwtSign"), SecretModel.JWTSign)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("body"), request.Body)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("bucketName"), os.Getenv("BucketName"))
	//Chequeo conexion a la base de datos o conecto la base de datos
	err = bd.ConectarBD(awsgo.Ctx)
	if err != nil {
		res = &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error conectando la base de datos " + err.Error(),
			Headers: map[string]string{
				"Content-Type": "Application/Json",
			},
		}
		return res, nil
	}

	respAPI := handlers.Manejadores(awsgo.Ctx, request)
	//Tengo que armar un APIGatewayProxyResponse. Si no me vino armado
	//un CustomRep desde las rutas, lo tengo que crear.
	//Cuando procese imagenes ya me va a venir un customResp totalmente
	//personalizado para esa ruta. Pero en las rutas donde venga sin
	//personalizar, voy a tener que armarme una respuesta
	if respAPI.CustomResp == nil {
		res = &events.APIGatewayProxyResponse{
			StatusCode: respAPI.Status,
			Body:       respAPI.Message,
			Headers: map[string]string{
				"Content-Type": "Application/Json",
			},
		}
		return res, nil
	} else {
		return respAPI.CustomResp, nil
	}
}

// os me va a permitir leer variables de entonro
// Es obligatorio que mi lambda reciba tres variables de entorno sino no funciona
func ValidoParametros() bool {
	_, traeParametro := os.LookupEnv("SecretName")
	if !traeParametro {
		return traeParametro
	}
	_, traeParametro = os.LookupEnv("BucketName")
	if !traeParametro {
		return traeParametro
	}
	_, traeParametro = os.LookupEnv("UrlPrefix")
	if !traeParametro {
		return traeParametro
	}
	return traeParametro
}
