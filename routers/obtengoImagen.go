package routers

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/LuisAli22/twitterGo/awsgo"
	"github.com/LuisAli22/twitterGo/bd"
	"github.com/LuisAli22/twitterGo/models"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

func ObtenerImagen(ctx context.Context, uploadType string, request events.APIGatewayProxyRequest, claim models.Claim) models.ResAPI {
	var r models.ResAPI
	r.Status = 400

	ID := request.QueryStringParameters["id"]
	if len(ID) < 1 {
		r.Message = "El paràmetro ID es obligatorio"
		return r
	}

	perfil, err := bd.BuscoPerfil(ID)
	if err != nil {
		r.Message = "Usuario no encontrado " + err.Error()
		return r
	}

	var filename string
	switch uploadType {
	case "A":
		filename = "avatars/" + perfil.Avatar
	case "B":
		filename = "banners/" + perfil.Banner
	}
	fmt.Println("Filename " + filename)
	svc := s3.NewFromConfig(awsgo.Cfg)
	file, err := downloadFromS3(ctx, svc, filename)
	if err != nil {
		r.Status = 500
		r.Message = "Error descargando archivo de S3 " + err.Error()
		return r
	}
	r.CustomResp = &events.APIGatewayProxyResponse{
		StatusCode: 200,
		//Notar que el file que es un archivo binario (un archivo de imagen)
		//puedo usar la funcion String() para convertirlo
		//Cada caracter, por mas binario que sea va a tener una representacion
		//de caracter dentro de mi string
		Body: file.String(),
		Headers: map[string]string{
			"Content-Type":        "application/octet-stream",
			"Content-Disposition": fmt.Sprintf("attachment; filename=\"%s\"", filename),
		},
	}
	return r
}

func downloadFromS3(ctx context.Context, svc *s3.Client, filename string) (*bytes.Buffer, error) {
	bucket := ctx.Value(models.Key("bucketName")).(string)
	obj, err := svc.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})
	if err != nil {
		return nil, err
	}
	//No es el body del request de postman. Se refiere
	//al body del objeto de s3
	//Lo abrimos con el GEtObject y despues hay que cerrarlo
	//por eso esta puesto en el defer, asi cuando salgamos
	//de la funcion, cierra el objeto
	defer obj.Body.Close()
	fmt.Println("bucketname = " + bucket)
	//Tengo que ir transformano el obj hasta tener un
	//formato para poder enviarlo al front
	file, err := io.ReadAll(obj.Body)
	if err != nil {
		return nil, err
	}
	buffer := bytes.NewBuffer(file)
	return buffer, nil

}
