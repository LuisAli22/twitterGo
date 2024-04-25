package routers

import (
	"bytes"
	"context"
	"encoding/base64"
	"io"
	"mime"
	"mime/multipart"
	"strings"

	"github.com/LuisAli22/twitterGo/bd"
	"github.com/LuisAli22/twitterGo/models"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type readSeeker struct {
	io.Reader
}

//En formdata vamos a enviar un archivo en formato multipart
//A diferencia de una rutina normal que toma un archivo del disco y lo sube a s3
// es que estamos corriendo en una lambda y la lambda de amazon es serverless.
//No estamos capturando la imagen de un archivo sino que la estamos recibiendo
//via API y la estamos subiendo en binario en multipart hacia s3.

// La siguiente funcion implementa una interfaz de tipo readSeeker (una estrucutra)
func (rs *readSeeker) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func UploadImage(ctx context.Context, uploadType string, request events.APIGatewayProxyRequest, claim models.Claim) models.ResAPI {
	var r models.ResAPI
	r.Status = 400
	IDUsuario := claim.Id.Hex()

	var filename string
	var usuario models.Usuario
	//Capturo el bucket de nuestro context
	bucket := aws.String(ctx.Value(models.Key("bucketName")).(string))
	switch uploadType {
	case "A":
		filename = "avatars/" + IDUsuario + ".jpg"
		usuario.Avatar = filename
	case "B":
		filename = "banners/" + IDUsuario + ".jpg"
		usuario.Banner = filename
	}

	mediaType, params, err := mime.ParseMediaType(request.Headers["Content-Type"])
	if err != nil {
		r.Status = 500
		r.Message = err.Error()
		return r
	}

	//Alguien podria mandar el archivo en formato binario y no me sirve
	//tiene que ser multipart
	if strings.HasPrefix(mediaType, "multipart/") {
		//La imagen viene en el body
		body, err := base64.StdEncoding.DecodeString(request.Body)
		if err != nil {
			//Hubo en error al decodificar
			r.Status = 500
			r.Message = err.Error()
			return r
		}
		//Lo va a crear en memoria como si lo hubiera leido del disco
		mr := multipart.NewReader(bytes.NewReader(body), params["boundary"])
		//Una vex creado, el Next se posiciona y te dice realmente si pudo convertir
		//la imagen en NewReader. Es como cuando en base de datos tenes un cursor y
		//le decis next para que se posiciones en el proximo registro.
		p, err := mr.NextPart()
		//cuando llega a EOF tambien provoca un error el Next pero
		//si es EOF esta bien que suceda. El tema es si pasa un error y no
		//es porque se hizo un Next cuando llego a EOF
		if err != nil && err != io.EOF {
			r.Status = 500
			r.Message = err.Error()
			return r
		}

		if err != io.EOF {
			if p.FileName() != "" {
				//Creo un buffer vacio
				//Cuando leemos un archivo desde disco, lo aloja en un buffer en memoria
				//temporal. Cuando el archivo se cierra, el area de moemoria queda eliminada
				buf := bytes.NewBuffer(nil)
				//Copio el archivo al buffer
				if _, err := io.Copy(buf, p); err != nil {
					//Hubo un error al copiar
					r.Status = 500
					r.Message = err.Error()
					return r
				}
				//Creo una nueva sesion de amazon. Configuro la region
				//para la region de virginia (us-east-1)
				sess, err := session.NewSession(&aws.Config{
					Region: aws.String("us-east-1")})
				if err != nil {
					r.Status = 500
					r.Message = err.Error()
					return r
				}
				//Abor el camino hacia el bucket s3
				uploader := s3manager.NewUploader(sess)
				//Ahora creamos el canal para poder subir el archivo a s3
				//En el UplaodInput van todos los parametros para poder subir
				_, err = uploader.Upload(&s3manager.UploadInput{
					Bucket: bucket,
					Key:    aws.String(filename),
					Body:   &readSeeker{buf},
				})
				if err != nil {
					r.Status = 500
					r.Message = err.Error()
					return r
				}
			}
		}
		//Llamar al paquet de base de datos para que grabe en mongo
		status, err := bd.ModificoRegistro(usuario, IDUsuario)
		if err != nil || !status {
			r.Status = 400
			r.Message = "Error al modificar registro del usuario " + err.Error()
			return r
		}
		r.Status = 200
		r.Message = "Image Uplaoded OK!"
		return r

	} else {
		r.Message = "Debe enviar una imagen con el 'Context-ype' de tipo `multipart` en el HEader"
		r.Status = 400
		return r
	}

}
