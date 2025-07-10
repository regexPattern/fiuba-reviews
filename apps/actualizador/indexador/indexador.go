package indexador

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type Indexador struct {
	client     *s3.Client
	bucketName *string
}

type Oferta struct {
	Materias []Materia
}

type Materia struct {
	Nombre   string    `json:"nombre"`
	Codigo   string    `json:"codigo"`
	Catedras []Catedra `json:"catedras"`
}

type Catedra struct {
	Codigo   int       `json:"codigo"`
	Docentes []Docente `json:"docentes"`
}

type Docente struct {
	Nombre string `json:"nombre"`
	Rol    string `json:"rol"`
}

func New(client *s3.Client, bucketName *string) *Indexador {
	return &Indexador{client: client, bucketName: bucketName}
}

func (i *Indexador) IndexarOfertasComisiones() ([]*Oferta, error) {
	objs, _ := i.getObjectsFromBucket()
	var ofertas []*Oferta

	for _, o := range objs {
		o, _ := i.newOferta(o)
		ofertas = append(ofertas, o)
	}

	return ofertas, nil
}

func (i *Indexador) getObjectsFromBucket() ([]s3types.Object, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	output, err := i.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: i.bucketName,
	})
	if err != nil {
		return nil, err
	}

	slog.Info(fmt.Sprintf("obtenidos %v archivos del bucket", len(output.Contents)))

	return output.Contents, nil
}

func (i *Indexador) newOferta(obj s3types.Object) (*Oferta, error) {
	return nil, nil
}

func (i *Indexador) newMateria(obj s3types.Object) (*Materia, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	head, _ := i.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: i.bucketName,
		Key:    obj.Key,
	})

	fmt.Println(head)

	return nil, nil
}
