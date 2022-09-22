package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func respuesta(response http.ResponseWriter, res []byte, err error) {
	if err != nil {
		respuestaError(response, http.StatusBadRequest, "error en la consulta de base de datos")
		return
	}
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(200)
	response.Write(res)

}

func respuestaError(response http.ResponseWriter, estado int, mensaje string) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(estado)
	response.Write([]byte(`{"error":"` + mensaje + `"}`))
}

func LevantarServidor() {
	router := mux.NewRouter().StrictSlash(true)
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "PUT"})
	origins := handlers.AllowedOrigins([]string{"*"})
	port := os.Getenv("PORT")
	if port == "" {
		port = "10000"
	}
	router.HandleFunc("/", inicio).Methods("GET")

	router.HandleFunc("/getAllUsers", AllUsers).Methods("GET")
	router.HandleFunc("/getAllSessions", Sessions).Methods("GET")
	router.HandleFunc("/getAllPets", AllPets).Methods("GET")
	router.HandleFunc("/getAllLogs", AllLogs).Methods("GET")
	router.HandleFunc("/User/{id}", UsuarioId).Methods("GET")
	router.HandleFunc("/Pet/{id}", PetId).Methods("GET")
	router.HandleFunc("/Session", SessionId).Methods("POST")
	router.HandleFunc("/UploadPhoto", UploadPhoto).Methods("POST") //PENDIENTE
	router.HandleFunc("/AddUser", AddUser).Methods("POST")
	router.HandleFunc("/AddPet", AddPet).Methods("POST")
	router.HandleFunc("/AddSession", AddSession).Methods("POST")
	router.HandleFunc("/AddSessionImage", AddSessionImage).Methods("POST")
	router.HandleFunc("/AddSessionMed", AddSessionMed).Methods("POST")
	router.HandleFunc("/AddLog", AddLog).Methods("POST")
	router.HandleFunc("/AcceptUser/{id}", AcceptUser).Methods("GET")
	router.HandleFunc("/AcceptSession", AcceptSession).Methods("POST")
	router.HandleFunc("/UpdatePrice", UpdatePrice).Methods("POST")
	router.HandleFunc("/GetPrice", GetTotalPrice).Methods("POST")

	fmt.Println("puerto " + port)
	http.ListenAndServe(":"+port, handlers.CORS(headers, methods, origins)(router))
}

func inicio(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(200)
	response.Write([]byte(`{"API":"welcome"}`))
}

/********************* FUNCIONES GET PARA OBTENER TODOS LOS REGISTROS *********************/
func AllUsers(response http.ResponseWriter, request *http.Request) {
	users := getAllUsers()
	res, err := json.Marshal(users)
	respuesta(response, res, err)
}

func Sessions(response http.ResponseWriter, request *http.Request) {
	sessions := getAllSessions()
	res, err := json.Marshal(sessions)
	respuesta(response, res, err)
}

func AllPets(response http.ResponseWriter, request *http.Request) {
	pets := getAllPets()
	res, err := json.Marshal(pets)
	respuesta(response, res, err)
}

func AllLogs(response http.ResponseWriter, request *http.Request) {
	logs := getAllLogs()
	res, err := json.Marshal(logs)
	respuesta(response, res, err)
}

/********************* FUNCIONES GET PARA OBTENER ELEMENTO POR ID *********************/
func UsuarioId(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["id"]
	user := getUsuario(id)
	res, err := json.Marshal(user)
	respuesta(response, res, err)

}

func PetId(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["id"]
	user := getMascota(id)
	res, err := json.Marshal(user)
	respuesta(response, res, err)
}

func SessionId(response http.ResponseWriter, request *http.Request) {
	data, err := ioutil.ReadAll(request.Body)
	if err != nil {
		respuestaError(response, http.StatusBadRequest, "error en el envio de datos")
		return
	}
	// en el json se recibe id
	var SessionImage map[string]interface{}
	err2 := json.Unmarshal(data, &SessionImage)
	if err2 != nil {
		respuestaError(response, http.StatusBadRequest, "error al convertir datos")
		return
	}

	id := SessionImage["id"].(string)
	user := getSesion(id)
	res, err := json.Marshal(user)
	respuesta(response, res, err)
}

// Funcion para agregar foto al S3 PENDIENTE
func UploadPhoto(response http.ResponseWriter, request *http.Request) {

}

// Funcion para agregar usuario
func AddUser(response http.ResponseWriter, request *http.Request) {
	data, err := ioutil.ReadAll(request.Body)
	if err != nil {
		respuestaError(response, http.StatusBadRequest, "error en el envio de datos")
		return
	}
	var user Usuario
	err2 := json.Unmarshal(data, &user)
	if err2 != nil {
		respuestaError(response, http.StatusBadRequest, "error al convertir datos")
		return
	}
	if user.Tipo == "cliente" || user.Tipo == "secretaria" {
		user.Especialidad = user.Tipo
	}
	if !insertUsuario(user) {
		respuestaError(response, http.StatusBadRequest, "error al ingresar datos en la base de datos")
		return
	}
	respuesta(response, []byte(`{"response":"usuario agregado exitosamente"}`), nil)
}

// Funcion para agregar una session
func AddSession(response http.ResponseWriter, request *http.Request) {
	data, err := ioutil.ReadAll(request.Body)
	if err != nil {
		respuestaError(response, http.StatusBadRequest, "error en el envio de datos")
		return
	}
	var sesion Sesion
	err2 := json.Unmarshal(data, &sesion)
	if err2 != nil {
		respuestaError(response, http.StatusBadRequest, "error al convertir datos")
		return
	}
	if !insertSesion(sesion) {
		respuestaError(response, http.StatusBadRequest, "error al ingresar datos en la base de datos")
		return
	}
	respuesta(response, []byte(`{"response":"session agregada exitosamente"}`), nil)
}

// Funcion para agregar un log
func AddLog(response http.ResponseWriter, request *http.Request) {
	data, err := ioutil.ReadAll(request.Body)
	if err != nil {
		respuestaError(response, http.StatusBadRequest, "error en el envio de datos")
		return
	}
	var log Logs
	err2 := json.Unmarshal(data, &log)
	if err2 != nil {
		respuestaError(response, http.StatusBadRequest, "error al convertir datos")
		return
	}
	if !insertLog(log) {
		respuestaError(response, http.StatusBadRequest, "error al ingresar datos en la base de datos")
		return
	}
	respuesta(response, []byte(`{"response":"Log agregado exitosamente"}`), nil)
}

// Funcion para agregar mascota
func AddPet(response http.ResponseWriter, request *http.Request) {
	data, err := ioutil.ReadAll(request.Body)
	if err != nil {
		respuestaError(response, http.StatusBadRequest, "error en el envio de datos")
		return
	}
	var petProfile MascotaImage
	err2 := json.Unmarshal(data, &petProfile)
	if err2 != nil {
		respuestaError(response, http.StatusBadRequest, "error al convertir datos")
		return
	}
	var pet Mascota
	pet.IdMascota = petProfile.IdMascota
	pet.Nombre = petProfile.Nombre
	pet.Raza = petProfile.Raza
	pet.Edad = petProfile.Edad
	pet.Foto_Url = petProfile.Foto_Url
	pet.Username = petProfile.Username
	if !insertMascota(pet, petProfile.ProfilePhoto) {
		respuestaError(response, http.StatusBadRequest, "error al ingresar datos en la base de datos")
		return
	}
	respuesta(response, []byte(`{"response":"Mascota agregado exitosamente"}`), nil)
}

// Funcion para agregar fotos en sesion
func AddSessionImage(response http.ResponseWriter, request *http.Request) {
	data, err := ioutil.ReadAll(request.Body)
	if err != nil {
		respuestaError(response, http.StatusBadRequest, "error en el envio de datos")
		return
	}
	// en el json se recibe idSession y image
	var SessionImage map[string]interface{}
	err2 := json.Unmarshal(data, &SessionImage)
	if err2 != nil {
		respuestaError(response, http.StatusBadRequest, "error al convertir datos")
		return
	}
	if !addSessionImage(SessionImage["idSession"].(string), SessionImage["image"].(string)) {
		respuestaError(response, http.StatusBadRequest, "error al ingresar datos en la base de datos")
		return
	}
	respuesta(response, []byte(`{"response":"Log agregado exitosamente"}`), nil)
}

func AddSessionMed(response http.ResponseWriter, request *http.Request) {
	data, err := ioutil.ReadAll(request.Body)
	if err != nil {
		respuestaError(response, http.StatusBadRequest, "error en el envio de datos")
		return
	}
	// en el json se recibe idSession y med
	var SessionMed map[string]interface{}
	err2 := json.Unmarshal(data, &SessionMed)
	if err2 != nil {
		respuestaError(response, http.StatusBadRequest, "error al convertir datos")
		return
	}
	if !addMedicine(SessionMed["idSession"].(string), SessionMed["med"].(string)) {
		respuestaError(response, http.StatusBadRequest, "error al ingresar datos en la base de datos")
		return
	}
	respuesta(response, []byte(`{"response":"Log agregado exitosamente"}`), nil)
}

func AcceptUser(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["id"]
	err := confirmUser(id)
	if !err {
		respuestaError(response, http.StatusBadRequest, "error al ingresar datos en la base de datos")
		return
	}
	respuesta(response, []byte(`{"response":"usuario aceptado exitosamente"}`), nil)
}

func AcceptSession(response http.ResponseWriter, request *http.Request) {
	data, err := ioutil.ReadAll(request.Body)
	if err != nil {
		respuestaError(response, http.StatusBadRequest, "error en el envio de datos")
		return
	}
	// en el json se recibe idSession y med
	var inputSession map[string]interface{}
	err2 := json.Unmarshal(data, &inputSession)
	if err2 != nil {
		respuestaError(response, http.StatusBadRequest, "error al convertir datos")
		return
	}
	err3 := acceptSession(inputSession["id"].(string))
	if !err3 {
		respuestaError(response, http.StatusBadRequest, "error al ingresar datos en la base de datos")
		return
	}
	respuesta(response, []byte(`{"response":"usuario aceptado exitosamente"}`), nil)
}

func UpdatePrice(response http.ResponseWriter, request *http.Request) {
	data, err := ioutil.ReadAll(request.Body)
	if err != nil {
		respuestaError(response, http.StatusBadRequest, "error en el envio de datos")
		return
	}
	// en el json se recibe categoria y precio
	var precios map[string]interface{}
	err2 := json.Unmarshal(data, &precios)
	if err2 != nil {
		respuestaError(response, http.StatusBadRequest, "error al convertir datos")
		return
	}
	retorno := updatePrice(precios["categoria"].(string), precios["precio"].(float64))
	if !retorno {
		respuestaError(response, http.StatusBadRequest, "error al ingresar datos en la base de datos")
		return
	}
	respuesta(response, []byte(`{"response":"precio actualizado"}`), nil)
}

func GetTotalPrice(response http.ResponseWriter, request *http.Request) {
	data, err := ioutil.ReadAll(request.Body)
	if err != nil {
		respuestaError(response, http.StatusBadRequest, "error en el envio de datos")
		return
	}
	// en el json se recibe username, categoria y tiempo
	var precios map[string]interface{}
	err2 := json.Unmarshal(data, &precios)
	if err2 != nil {
		respuestaError(response, http.StatusBadRequest, "error al convertir datos")
		return
	}
	username := precios["username"].(string)
	categoria := precios["categoria"].(string)
	tiempo := precios["tiempo"].(float64)
	numeroVisitas := getCountVisits(username)
	subTotal := getPrice(categoria, tiempo)
	if numeroVisitas > 2 {
		subTotal = subTotal * 0.85
	}
	respuesta(response, []byte(`{"total":`+fmt.Sprintf("%f", subTotal)+`}`), nil)
}
