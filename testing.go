package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func testClients() {
	// List S3 Objects
	output, err := client_s3.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket_name),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("first page results:")
	for _, object := range output.Contents {
		log.Printf("key=%s size=%d", aws.ToString(object.Key), object.Size)
	}
	// List Dynamo Tables
	resp, err := client_dynamo.ListTables(context.TODO(), &dynamodb.ListTablesInput{
		Limit: aws.Int32(5),
	})
	if err != nil {
		log.Fatalf("failed to list tables, %v", err)
	}
	fmt.Println("Tables:")
	for _, tableName := range resp.TableNames {
		fmt.Println(tableName)
	}
}

func testFunctions(caso int, insert bool) {
	switch caso {
	case 1:
		if insert {
			insertUsuario(Usuario{Username: "ldecast", Nombre: "Luis Danniel Castellanos", Correo: "luis.danniel@hotmail.com", Contrasena: "password123", Tipo: "Cliente", Estado: 1})
		}
		u := getAllUsers()
		fmt.Println("Todos los usuarios:")
		for _, v := range u {
			fmt.Println(v)
			fmt.Println()
		}
		x := getUsuario("ldecast")
		fmt.Println("Usuario por id ldescast:")
		fmt.Println(x)
	case 2:
		if insert {
			insertMascota(Mascota{IdMascota: "M0001", Nombre: "Kira", Raza: "Mestiza", Edad: 3, Foto_Url: "ldecast/M0001/foto.jpg", Username: "ldecast"}, readFile("b64.txt"))
		}
		u := getAllPets()
		fmt.Println("Todos las mascotas:")
		for _, v := range u {
			fmt.Println(v)
			fmt.Println()
		}
		x := getMascota("M0001")
		fmt.Println("Mascota por id M0001:")
		fmt.Println(x)
	case 3:
		if insert {
			insertSesion(Sesion{IdSesion: "S0002", Creacion: time.Now().Format("02/01/2006 15:04:05"), Estado: 1, Hora_Ingreso: (time.Now().Add(time.Hour).Format("02/01/2006 15:04:05")), Hora_Salida: (time.Now().Add(time.Hour * 2).Format("02/01/2006 15:04:05")), Medicamentos: []string{"Galletas", "Anti pulgas", "Vitaminas"}, Tipo: "Cita", Imagenes: []string{}, IdMascota: "M0001"})
		}
		u := getAllSessions()
		fmt.Println("Todas las sesiones:")
		for _, v := range u {
			fmt.Println(v)
			fmt.Println()
		}
		x := getSesion("S0002")
		fmt.Println("Sesion por id S0002:")
		fmt.Println(x)
	case 4:
		if insert {
			addSessionImage("S0002", readFile("b64.txt"))
		}
		x := getSesion("S0002")
		fmt.Println("Imágenes en sesión S0002:")
		fmt.Println(x.Imagenes)
	case 5:
		if insert {
			m := make(map[string]string)
			m["campo1"] = "valor1"
			m["campo2"] = "valor2"
			insertLog(Logs{IdLog: "L0001", Metodo: "/cita", Entrada: m, Salida: m, Error: 0, Fecha_hora: time.Now().Format("02/01/2006 15:04:05"), Unix_Timestamp: int(time.Now().UTC().UnixNano())})
		}
		u := getAllLogs()
		fmt.Println("Todos los logs:")
		for _, v := range u {
			fmt.Println(v)
			fmt.Println()
		}
	}
}

func readFile(filename string) string {
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return string(content)
}
