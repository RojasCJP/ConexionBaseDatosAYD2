package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const bucket_name string = "gatifu-bucket"

var client_dynamo *dynamodb.Client
var client_s3 *s3.Client
var uploader_s3 *manager.Uploader

func instanceClients() {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-2"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	// Create an Amazon S3 service client and set uploader
	client_s3 = s3.NewFromConfig(cfg)
	uploader_s3 = manager.NewUploader(client_s3)
	// Create the DynamoDB client
	client_dynamo = dynamodb.NewFromConfig(cfg)
}

/********************* FUNCIONES PARA OBTENER TODOS LOS REGISTROS *********************/
func getAllUsers() []Usuario {
	out, err := client_dynamo.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String("Usuario"),
	})
	if err != nil {
		panic(err)
	}
	var users []Usuario
	err = attributevalue.UnmarshalListOfMaps(out.Items, &users)
	if err != nil {
		log.Printf("Couldn't unmarshal query response: %v\n", err)
	}
	return users
}

func getAllPets() []Mascota {
	out, err := client_dynamo.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String("Mascota"),
	})
	if err != nil {
		panic(err)
	}
	var pets []Mascota
	err = attributevalue.UnmarshalListOfMaps(out.Items, &pets)
	if err != nil {
		log.Printf("Couldn't unmarshal query response: %v\n", err)
	}
	return pets
}

func getAllSessions() []Sesion {
	out, err := client_dynamo.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String("Sesion"),
	})
	if err != nil {
		panic(err)
	}
	var sessions []Sesion
	err = attributevalue.UnmarshalListOfMaps(out.Items, &sessions)
	if err != nil {
		log.Printf("Couldn't unmarshal query response: %v\n", err)
	}
	return sessions
}

func getAllLogs() []Logs {
	out, err := client_dynamo.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String("Logs"),
	})
	if err != nil {
		panic(err)
	}
	var logs []Logs
	err = attributevalue.UnmarshalListOfMaps(out.Items, &logs)
	if err != nil {
		log.Printf("Couldn't unmarshal query response: %v\n", err)
	}
	return logs
}

/********************* FUNCIONES PARA OBTENER ELEMENTO POR ID *********************/
func getUsuario(idUsuario string) Usuario {
	out, err := client_dynamo.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("Usuario"),
		Key: map[string]types.AttributeValue{
			"Username": &types.AttributeValueMemberS{Value: idUsuario},
		},
	})
	if err != nil {
		log.Printf("Error returning item from Dynamo: %v\n", err)
	}
	var user Usuario
	err = attributevalue.UnmarshalMap(out.Item, &user)
	if err != nil {
		log.Printf("Couldn't unmarshal query response: %v\n", err)
	}
	return user
}

func getMascota(idMascota string) Mascota {
	out, err := client_dynamo.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("Mascota"),
		Key: map[string]types.AttributeValue{
			"IdMascota": &types.AttributeValueMemberS{Value: idMascota},
		},
	})
	if err != nil {
		log.Printf("Error returning item from Dynamo: %v\n", err)
	}
	var pet Mascota
	err = attributevalue.UnmarshalMap(out.Item, &pet)
	if err != nil {
		log.Printf("Couldn't unmarshal query response: %v\n", err)
	}
	return pet
}

func getSesion(idSesion string) Sesion {
	out, err := client_dynamo.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("Sesion"),
		Key: map[string]types.AttributeValue{
			"IdSesion": &types.AttributeValueMemberS{Value: idSesion},
		},
	})
	if err != nil {
		log.Printf("Error returning item from Dynamo: %v\n", err)
	}
	var session Sesion
	err = attributevalue.UnmarshalMap(out.Item, &session)
	if err != nil {
		log.Printf("Couldn't unmarshal query response: %v\n", err)
	}
	return session
}

/********************* FUNCIONES PARA GUARDAR UNA IMAGEN EN S3 *********************/
func uploadPhoto(base64_photo string, filepath string) bool {
	// Procesar imagen
	decode, err := base64.StdEncoding.DecodeString(base64_photo)
	if err != nil {
		log.Printf("Error decoding the base64 string: %v\n", err)
	}
	_, err = uploader_s3.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucket_name),
		Key:    aws.String(filepath),
		Body:   bytes.NewReader(decode),
	})
	if err != nil {
		log.Printf("Couldn't upload image: %v\n", err)
		return false
	} else {
		fmt.Println("Image " + filepath + " uploaded succesfully.")
		return true
	}
}

/********************* FUNCIONES PARA INSERTAR UN ELEMENTO *********************/
func insertUsuario(user Usuario) bool {
	item, err := attributevalue.MarshalMap(user)
	_, err = client_dynamo.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("Usuario"), Item: item,
	})
	if err != nil {
		log.Printf("Couldn't add item to table: %v\n", err)
		return false
	}
	return true
	// fmt.Println(out.Attributes)
}

func insertMascota(pet Mascota, base64_mascota string) bool {
	pet.Foto_Url = pet.Username + "/" + pet.IdMascota + "/foto.jpg"
	item, err := attributevalue.MarshalMap(pet)
	_, err = client_dynamo.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("Mascota"), Item: item,
	})
	if err != nil {
		log.Printf("Couldn't add item to table: %v\n", err)
		return false
	} else {
		uploadPhoto(base64_mascota, pet.Foto_Url)
	}
	return true
}

func insertSesion(session Sesion) bool {
	item, err := attributevalue.MarshalMap(session)
	_, err = client_dynamo.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("Sesion"), Item: item,
	})
	if err != nil {
		log.Printf("Couldn't add item to table: %v\n", err)
		return false
	}
	return true
}

func addSessionImage(id_session string, base64_image string) bool {
	session := getSesion(id_session)
	new_url := "sesiones/" + id_session + "/imagen" + strconv.Itoa(len(session.Imagenes)+1) + ".jpg"
	if uploadPhoto(base64_image, new_url) {
		session.Imagenes = append(session.Imagenes, new_url)
		_, err := client_dynamo.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
			TableName: aws.String("Sesion"),
			Key: map[string]types.AttributeValue{
				"IdSesion": &types.AttributeValueMemberS{Value: id_session},
			},
			UpdateExpression: aws.String("set Imagenes = :Imagenes"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":Imagenes": &types.AttributeValueMemberSS{Value: session.Imagenes},
			},
		})
		if err != nil {
			log.Printf("Couldn't add session image to array on Dynamo: %v\n", err)
			return false
		}
		// fmt.Println(out.Attributes)
	}
	return true
}

func insertLog(log_item Logs) bool {
	item, err := attributevalue.MarshalMap(log_item)
	_, err = client_dynamo.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("Logs"), Item: item,
	})
	if err != nil {
		log.Printf("Couldn't add item to table: %v\n", err)
		return false
	}
	return true
}

/********************* AGREGAR MEDICAMENTO A UNA SESION *********************/
func addMedicine(id_session string, medicamento string) bool {
	session := getSesion(id_session)
	session.Medicamentos = append(session.Medicamentos, medicamento)
	_, err := client_dynamo.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: aws.String("Sesion"),
		Key: map[string]types.AttributeValue{
			"IdSesion": &types.AttributeValueMemberS{Value: id_session},
		},
		UpdateExpression: aws.String("set Medicamentos = :Medicamentos"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":Medicamentos": &types.AttributeValueMemberSS{Value: session.Medicamentos},
		},
	})
	if err != nil {
		log.Printf("Couldn't add session medicine to array on Dynamo: %v\n", err)
	} else {
		fmt.Println("Medicamentos: ", session.Medicamentos)
	}
	return true
}

/********************* CONFIRMAR SESION (DE 0 PASA A 1) *********************/
func acceptSession(id_session string) bool {
	session := getSesion(id_session)
	session.Estado = 1
	_, err := client_dynamo.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: aws.String("Sesion"),
		Key: map[string]types.AttributeValue{
			"IdSesion": &types.AttributeValueMemberS{Value: id_session},
		},
		UpdateExpression: aws.String("set Estado = :Estado"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":Estado": &types.AttributeValueMemberN{Value: "1"},
		},
	})
	if err != nil {
		log.Printf("Couldn't add session medicine to array on Dynamo: %v\n", err)
	} else {
		fmt.Println("Medicamentos: ", session.Medicamentos)
	}
	return true
}

/********************* CONFIRMAR USUARIO (DE 0 PASA A 1) *********************/
func confirmUser(idUsuario string) bool {
	_, err := client_dynamo.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: aws.String("Usuario"),
		Key: map[string]types.AttributeValue{
			"Username": &types.AttributeValueMemberS{Value: idUsuario},
		},
		UpdateExpression: aws.String("set Estado = :Estado"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":Estado": &types.AttributeValueMemberN{Value: "1"},
		},
	})
	if err != nil {
		log.Printf("Couldn't add session medicine to array on Dynamo: %v\n", err)
		return false
	}
	log.Printf("Usuario confirmado: %v\n", idUsuario)
	return true
}

/********************* FUNCION PARA OBTENER EL MONTO A COBRAR *********************/
func getPrice(motivo string, horas float64) float64 {
	out, err := client_dynamo.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("Precios"),
		Key: map[string]types.AttributeValue{
			"Motivo": &types.AttributeValueMemberS{Value: motivo},
		},
	})
	if err != nil {
		log.Printf("Error returning price item from Dynamo: %v\n", err)
	}
	var precio Precios
	err = attributevalue.UnmarshalMap(out.Item, &precio)
	if err != nil {
		log.Printf("Couldn't unmarshal query response: %v\n", err)
	}
	monto := precio.Precio * horas
	fmt.Println("El monto a cobrar es:", monto)
	return monto
}

/********************* CONTAR CANTIDAD DE VISITAS (PARA SABER SI ES FRECUENTE) *********************/
func getCountVisits(idUsuario string) int {
	sessions := getAllSessions()
	pets := getAllPets()
	var count = 0
	for _, session := range sessions {
		for _, pet := range pets {
			if session.IdMascota == pet.IdMascota {
				if pet.Username == idUsuario {
					count++
					break
				}
			}
		}
	}
	fmt.Println("Visitas del cliente:", count)
	return count
}

/********************* MODIFICAR VALORES DE LA TABLA PRECIOS *********************/
func updatePrice(nombre string, nuevoValor float64) bool {
	_, err := client_dynamo.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: aws.String("Precios"),
		Key: map[string]types.AttributeValue{
			"Motivo": &types.AttributeValueMemberS{Value: nombre},
		},
		UpdateExpression: aws.String("set Precio = :Precio"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":Precio": &types.AttributeValueMemberN{Value: fmt.Sprintf("%v", nuevoValor)},
		},
	})
	if err != nil {
		log.Printf("Couldn't update price: %v\n", err)
		return false
	}
	return true
}

func main() {
	instanceClients()
	testClients()
	LevantarServidor()
	// testFunctions(1, false)
}
